package api

import "github.com/gin-gonic/gin"

func SetRouter(e *gin.Engine) {
	e.POST("/login", Login)
	e.POST("/register", Register)
	//弹幕接口
	e.GET("/ws", WsHandler)
	//红包生成接口  相当于发红包

	g := e.Group("/red")
	{
		g.POST("/creat", CreatRedPacket)
		//抢红包
		g.POST("/consume", ConsumeRedPacket)
	}
	//创建抽奖
	e.POST("/creatlottery", CreatLottery)
	//添加黑名单
	e.POST("/addblack", AddBlackList)
}
