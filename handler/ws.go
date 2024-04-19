package handler

import (
	"garptea/ws"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

func WsHandler(c *websocket.Conn) {

	claims := c.Locals("claims").(*casdoorsdk.Claims)
	connectionParams := ws.ConnectionParams{
		UserId:     claims.User.Id,
		Connection: c,
	}

	defer ws.Unregister(connectionParams)

	ws.Register(connectionParams)

	c.SetPongHandler(func(appData string) error {
		log.Info("pong")
		// if err := auth.ValidateUserAndSession(claims); err != nil {
		// 	ws.Unregister(connectionParams)
		// }
		return nil
	})
	c.SetCloseHandler(func(code int, text string) error {
		return nil
	})
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Errorf("read error: %s, user_id=%s", err, connectionParams.UserId)
			}
			return
		}
		if messageType == websocket.TextMessage {
			log.Info(message)
			ws.ProcessMessage(ws.ClientMessage{
				Connection: c,
				Message:    string(message),
			})
		}
	}
}
