package timer

import (
	"exam/limiter"
	"time"
)

func StartTimer() {
	t := time.Now()
	//获取时间
	next := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()+10, 0, 0, t.Location())
	//设置timer
	timer := time.NewTimer(next.Sub(t))
	<-timer.C
	limiter.L.ReSetLimiter()
	//无限循环
	StartTimer()
}
