package server

import (
	"bytes"
	"encoding/json"
	"exam/limiter"
	"exam/tree"
	"fmt"
	"github.com/gorilla/websocket"
	"image/color"
	"log"
	"strings"
	"time"
)

var M *Manger

var ltr *Lottery

type Manger struct {
	Clients []*Client
	Lot     chan *Lottery //抽奖
	Over    chan int      //时间过了
	Ing     bool
}

type Client struct {
	Name    string          `json:"name"`   //自己的名字
	Socket  *websocket.Conn `json:"socket"` //stock里面的东西
	SendMsg chan []byte     `json:"send"`   //[]byte类型的msg，更便于接收
	In      chan int        `json:"in"`     //进入直播间
	Out     chan int
	IsBlack bool
}

func InitManger() {
	M = &Manger{
		Clients: make([]*Client, 10000),
		Lot:     make(chan *Lottery),
		Over:    make(chan int),
	}
	go findLottery()
}

//构造函数
func NewClient(name string, conn *websocket.Conn, b bool) *Client {
	c := &Client{
		Name:    name,
		Socket:  conn,
		SendMsg: make(chan []byte, 1000),
		In:      make(chan int, 1),
		Out:     make(chan int, 1),
		IsBlack: b,
	}
	M.Clients = append(M.Clients, c)
	return c
}

type msgType uint8

const (
	ordinary msgType = iota //普通的弹幕
	push                    //推送给所有人 比如：谁谁中奖了
	welcome
)

type Msg struct {
	Name    string     `json:"name"`
	Color   color.RGBA `json:"color"`   //rgb  颜色
	Time    int64      `json:"time"`    //假设传过来的是时间戳
	Content string     `json:"content"` //内容
	Code    msgType    `json:"code"`    //推送消息的类型 1
}

//发送给其他人
func (c *Client) send(msg []byte, code msgType) {
	switch code {
	case ordinary:
		for _, v := range M.Clients {
			if v != nil && v.Name != c.Name {
				v.SendMsg <- msg
			}
		}
	default:
		for _, v := range M.Clients {
			if v != nil {
				v.SendMsg <- msg
			}
		}

	}

}

//删除退出的人
func (c *Client) deleteClient() {
	for k, v := range M.Clients {
		if v!=nil&&v.Name == c.Name {
			if k != len(M.Clients)-1 {
				M.Clients = append(M.Clients[:k], M.Clients[k+1:]...)
				return
			}
			M.Clients = M.Clients[:k]
		}
	}
}

//一个用户的开始
func (c *Client) OneStart() {
	for {
		select {
		case <-c.In:
			str := "!!!!!欢迎 <" + c.Name + ">进入房间"
			fmt.Println(str)
			jsonMsg := CreatBytesMsg("", color.RGBA{}, time.Now().Unix(), str, welcome)
			c.send(jsonMsg, welcome)
		case <-c.Out:
			close(c.SendMsg)
			c.deleteClient()
		}
	}
}

//读取
func (c *Client) Read() {
	defer func() {
		c.Out <- 1
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage() //读取信息
		if err != nil {
			break
		}
		if len(limiter.L.Conn) >= 1000 {
			continue
		}
		//直接传过去
		var m Msg
		decoder := json.NewDecoder(bytes.NewReader(message))
		if err := decoder.Decode(&m); err != nil {
			log.Println("json message error")
			continue
		}

		//如果有抽奖
		if ltr != nil && ltr.Num != 0 {
			if c.isLucky(m.Content) {
				ltr.Num--
			}
		}

		//部分特殊符号无法处理
		content := strings.ReplaceAll(m.Content, " ", "") //去除空格
		//判断是否有敏感词
		if str, b := tree.T.IsContain(content); b {
			m.Content = strings.ReplaceAll(content, str, "*")
			message, err = json.Marshal(m)
			if err != nil {
				log.Println("json error ", err)
			}

		}
		//fmt.Println(m.Content)

		limiter.L.SetLimiter()
		c.send(message, m.Code)

	}
}

//写入
func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendMsg:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte("get message error"))
				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message) //这里是输出message,以text
		}
	}
}

func CreatBytesMsg(name string, col color.RGBA, t int64, content string, code msgType) []byte {
	m := Msg{
		Name:    name,
		Color:   col,
		Time:    t,
		Content: content,
		Code:    code,
	}
	jsonMsg, _ := json.Marshal(m)
	return jsonMsg
}
