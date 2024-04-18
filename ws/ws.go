package ws

import (
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

var (
	clients    = make(map[*websocket.Conn]*client)
	register   = make(chan ConnectionParams)
	unregister = make(chan *websocket.Conn)
	command    = make(chan string)
)

type ConnectionParams struct {
	Connection *websocket.Conn
	UserId     string
}

type client struct {
	userId   string
	hbTicker *time.Ticker
}

func RunWsLoop() {
	for {
		select {
		case params := <-register:
			clients[params.Connection] = &client{
				userId:   params.UserId,
				hbTicker: time.NewTicker(time.Second * 5),
			}
			go runHeartBeatLoop(params.Connection)

		case connection := <-unregister:
			//client := clients[connection]

			delete(clients, connection)

		case cmd := <-command:
			log.Info("Command: %s", cmd)

		}
	}
}

func Unregister(c *websocket.Conn) {
	unregister <- c
	c.Close()
}

func Register(c ConnectionParams) {
	register <- c
}

func runHeartBeatLoop(c *websocket.Conn) {
	for {

	}
}
