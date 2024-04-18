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
	command    = make(chan ClientMessage)
)

type ConnectionParams struct {
	Connection *websocket.Conn
	UserId     string
}

type ClientMessage struct {
	Connection *websocket.Conn
	Message    string
}

type client struct {
	userId   string
	hbTicker *time.Ticker
}

func RunWsLoop() {
	for {
		select {
		case params := <-register:
			client := &client{
				userId:   params.UserId,
				hbTicker: time.NewTicker(time.Second * 5),
			}
			clients[params.Connection] = client
			go runHeartBeatLoop(params.Connection, client)

		case connection := <-unregister:
			cl := clients[connection]
			log.Info("Disconnect client user_id=%s", cl.userId)
			delete(clients, connection)

		case cmd := <-command:
			client := clients[cmd.Connection]
			log.Info("Command: %s", cmd.Message)
			client.hbTicker.Reset(time.Second * 5)
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

func ProcessMessage(msg ClientMessage) {
	command <- msg
}

func runHeartBeatLoop(c *websocket.Conn, client *client) {
	defer client.hbTicker.Stop()
	for {
		select {
		case <-client.hbTicker.C:
			err := c.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*3))
			if err != nil {
				log.Info("Disconnect client user_id=%s (no ping response)", client.userId)
				Unregister(c)
				return
			}
		case connection := <-unregister:
			if connection == c {
				return
			}
		}
	}
}
