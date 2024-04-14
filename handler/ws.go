package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type ClientData struct {
	Id         string
	Connection *websocket.Conn
}

type client struct {
}

var (
	clients    = make(map[*websocket.Conn]*client)
	register   = make(chan ClientData)
	broadcast  = make(chan string)
	unregister = make(chan *websocket.Conn)
)

func WsHandler(c *websocket.Conn) {
	c.SetPongHandler(func(msg string) error {
		fmt.Println("pong")
		return nil
	})

	go func() {
		time.Sleep(time.Second)
		err := c.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(time.Second*20))
		if err != nil {
			fmt.Println("err")
		}
	}()

	for {

	}
}

func RunWsLoop() {
	// for {
	// 	select {
	// 	case connection := <-register:
	// 		clients[connection] = &client{}
	// 	}
	// }
}
