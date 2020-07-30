package api

import (
	"exam/db"
	"exam/red_packet"
	"exam/rsp"
	"exam/server"
	"exam/tree"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"image/color"
	"log"
	"net/http"
	"strconv"
	"time"
)

const minRedPacket = 0.1

var (
	upgrade = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

//登录
func Login(c *gin.Context) {
	var u db.User
	if err := c.Bind(&u); err != nil {
		rsp.BindError(c, "bind user error")
		return
	}
	if db.FindUser(u.Name) == nil || db.FindUser(u.Name).PassWord != u.PassWord {
		rsp.LoginError(c, "password error or name error")
		return
	}
	rsp.Ok(c, "login success")
}

//注册
func Register(c *gin.Context) {
	var u db.User
	if err := c.Bind(&u); err != nil {
		rsp.BindError(c, "bind user error")
		return
	}
	if err := db.AddUser(u); err != nil {
		rsp.RegisterError(c, "register error")
		return
	}
	rsp.Ok(c, "register success")
}

//弹幕
func WsHandler(c *gin.Context) {
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		//升级失败 会话失败
		log.Println("websocket connect is failed,err:", err)

	}
	//从cookie中获取值
	name, err := c.Cookie("name")
	if err != nil {
		log.Println("get name error:", err)
		//return
		name = "as"
	}
	//是否是黑名单
	_, b := tree.NameTree.FindStr(name)
	client := server.NewClient(name, conn, b)
	client.In <- 1
	go client.OneStart()
	if !b {
		go client.Read() //只有是非黑名单，才能允许读取
	}
	go client.Write() //打印都允许
	//只要主线程不鸽，，协程都能跑
}

//请求生成红包
func CreatRedPacket(c *gin.Context) {
	var p red_packet.Packet
	if err := c.BindJSON(&p); err != nil {
		rsp.BindError(c, "json red packet error")
		return
	}
	reds := p.CreatRed(minRedPacket)

	rsp.Ok(c, "creat red  packet success")
	//生成一个红包
	red_packet.RedPackets[p.Id] = red_packet.NewPackets(p.Number)
	for _, v := range reds {
		red_packet.RedPackets[p.Id].OnePacket <- v
	}
	return
}

//请求抢红包
func ConsumeRedPacket(c *gin.Context) {
	name, err := c.Cookie("name")
	if err != nil || len(name) == 0 {
		rsp.GetDataError(c, "get name error")
		return
	}
	id, b := c.GetQuery("id")
	if b {
		rsp.GetDataError(c, "get red_id error")
		return
	}
	i, err := strconv.Atoi(id)
	if err != nil {
		rsp.OtherError(c, "transform id error")
		return
	}
	p, ok := red_packet.RedPackets[i]
	//过期了
	if !ok {
		rsp.NoRedPacket(c, "the red packet has expired")
		return
	}
	//没了
	if len(p.OnePacket) == 0 {
		rsp.NoRedPacket(c, "the red envelopes have been looted")
		delete(red_packet.RedPackets, i) //及时删除，不浪费储存空间
		return
	}
	//抢到了
	m := <-red_packet.RedPackets[i].OnePacket
	str := fmt.Sprintf("%s get %.2f yuan", name, m)
	//可以前端转一下
	rsp.GetRedPacket(c, str)
	//后端也可写
	for _, v := range server.M.Clients {
		jsonMsg := server.CreatBytesMsg("", color.RGBA{}, time.Now().Unix(), str, 1)
		v.SendMsg <- jsonMsg
	}
	return
}

//生成随机抽奖
func CreatLottery(c *gin.Context) {
	//用户验证
	//_, err := c.Cookie("name")
	//if err != nil {
	//	rsp.GetDataError(c, "get name error")
	//	return
	//}

	if server.M.Ing {
		rsp.PostError(c, "more lottery cannot be held")
		return
	}
	var ltr server.Lottery
	if err := c.BindJSON(&ltr); err != nil {
		rsp.BindError(c, "json lottery error")
		return
	}
	//抽奖持续时间至少大于等于3分钟
	if ltr.Minute < 3 {
		rsp.PostTimeToShort(c, "post time too short")
		return
	}
	rsp.Ok(c, "creat lottery success")
	go ltr.StartLottery()
	return

}

func AddBlackList(c *gin.Context) {
	var user db.BlackList
	if err := c.Bind(&user); err != nil {
		rsp.BindError(c, "json black list error")
		return
	}
	if err := db.AddBlack(user); err != nil {
		rsp.DbError(c, "add black user to db error")
		return
	}
	rsp.Ok(c, "add black success")
}
