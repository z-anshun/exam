package main

import (
	"encoding/json"
	"exam/server"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"image/color"
	"net/url"
	"time"
)
//101.201.140.26
var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	s := server.Msg{
		Name:    "as",
		Color:   color.RGBA{1, 2, 3, 1},
		Time:    time.Now().Unix(),
		Content: "test",
	}
	jsonMsg, _ := json.Marshal(s)
	conn.WriteMessage(websocket.TextMessage, []byte(jsonMsg))
	//defer  conn.Close()
	//循环读取，读得耍
	for {
		_, message, _ := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)

		}
		if len(message) != 0 {
			fmt.Printf("received: %s\n", message)
		}
	}
}
