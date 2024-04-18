package ws

import (
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

var (
	clients      = make(map[*websocket.Conn]*client)
	register     = make(chan ConnectionParams)
	unregister   = make(chan ConnectionParams)
	perUserMsg   = make(chan ClientMessage)
	broadcastMsg = make(chan ClientMessage)
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
			delete(clients, connection.Connection)

		case cmd := <-perUserMsg:
			cl := clients[cmd.Connection]
			log.Infof("Command: %s", cmd.Message)
			cl.hbTicker.Reset(time.Second * 5)

			//case msg := <-broadcastMsg:

		}
	}
}

func Unregister(c ConnectionParams) {
	cl := clients[c.Connection]
	log.Infof("Disconnect client user_id=%s", cl.userId)
	cl.hbTicker.Stop()
	c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
	c.Connection.Close()
	unregister <- c
}

func Register(c ConnectionParams) {
	register <- c
}

func ProcessMessage(msg ClientMessage) {
	perUserMsg <- msg
}

func runHeartBeatLoop(c *websocket.Conn, cl *client) {
	for {
		select {
		case <-cl.hbTicker.C:
			err := c.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*3))
			if err != nil {
				return
			}
		case connection := <-unregister:
			if cl.userId == connection.UserId {
				return
			}
		}
	}
}
