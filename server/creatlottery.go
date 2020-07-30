package server

import (
	"fmt"
	"time"
)

type Lottery struct {
	//抽奖个数
	Num int `json:"num"`
	//抽奖持续时间 分钟
	Minute int `json:"minute"`
	//指定弹幕内容
	Content string `json:"content"`
}

func (l *Lottery) StartLottery() {
	LuckyList = nil
	M.Ing = true
	go setTime(l.Minute)
	once := l.Num / 3
	M.Lot <- &Lottery{
		Num: once,
		Content:l.Content,
	}
	//开始的时候抽一下  中间时刻抽一下  结束前一分钟抽一下
	t := time.Now()
	//获取时间
	next := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()+(l.Minute/2-1), 0, 0, t.Location())
	//设置timer
	timer := time.NewTimer(next.Sub(t))
	//如果啥子商品都没用，就读取
	<-timer.C

	//中间时刻
	M.Lot <- &Lottery{
		Num: once,
		Content:l.Content,
	}
	//最后时刻
	timer.Reset(next.Sub(t))
	<-timer.C
	M.Lot <- &Lottery{
		Num: l.Num - once*2,
		Content:l.Content,
	}
}

//计时
func setTime(minute int) {
	//开始的时候抽一下  中间时刻抽一下  结束前一分钟抽一下
	t := time.Now()
	//获取时间
	next := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()+minute, t.Second(), t.Nanosecond(), t.Location())
	timer := time.NewTimer(next.Sub(t))
	<-timer.C
	fmt.Println("??")
	M.Over <- 1
}
