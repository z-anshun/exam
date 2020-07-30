package main

import (
	"exam/api"
	"github.com/gin-gonic/gin"
)

//{"name":"as","color":{"R":0,"G":0,"B":0,"A":1},"time":1596010334,"content":"no way"}
func main() {
	InitAll()
	e := gin.Default()
	api.SetRouter(e)

	e.Run(":8080")

}
