package rsp

import "github.com/gin-gonic/gin"

//0 -> 错误  1 -> 正确  2 -> 红包  3 -> 登录或注册错误 5 -> 数据库错误
//正常
func Ok(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": "100",
		"msg":  msg,
	})
}

//抢到红包
func GetRedPacket(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": "200",
		"msg":  msg,
	})
}

//未抢到红包
func NoRedPacket(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": "201",
		"msg":  msg,
	})
}

//登录错误
func LoginError(c *gin.Context,msg string){
	c.JSON(200, gin.H{
		"code": "301",
		"msg":  msg,
	})
}
//注册错误
func RegisterError(c *gin.Context,msg string){
	c.JSON(200, gin.H{
		"code": "302",
		"msg":  msg,
	})
}

//错误，虽然获取到了数据，但是有问题
func BindError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": "001",
		"msg":  msg,
	})
}

//错误 未获取到数据
func GetDataError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": "002",
		"msg":  msg,
	})
}

//请求错误
func PostError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": "003",
		"msg":  msg,
	})
}

//请求抽奖时间过短
func PostTimeToShort(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": "004",
		"msg":  msg,
	})
}

//其它未知错误
func OtherError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": "005",
		"msg":  msg,
	})
}

//数据库错误
func DbError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": "500",
		"msg":  msg,
	})
}
