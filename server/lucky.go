package server

import (
	"fmt"
	"image/color"
	"strings"
	"time"
)

//中奖人的名单
var LuckyList []string

//是否开始抽奖
func findLottery() {
	for {
		select {
		case l := <-M.Lot:

			if ltr != nil {
				ltr.Num += l.Num //分时间段抽取的，可能上一个时间段未抽完
			} else {
				ltr = l
			}

		case <-M.Over:
			ltr = nil
			fmt.Println("抽奖结束")
			M.Ing = false
			str := "恭喜：" + strings.Join(LuckyList, ",") + "中奖"
			msg := CreatBytesMsg("", color.RGBA{}, time.Now().Unix(), str, push)
			for _, v := range M.Clients {
				v.SendMsg <- msg
			}
		}
	}
}

func (c *Client) isLucky(str string) bool {
	//这货中奖了
	fmt.Println(str, ltr.Content)
	if str == ltr.Content {
		//写入名单
		fmt.Println(c.Name)
		LuckyList = append(LuckyList, c.Name)
		return true
	}
	return false
}
