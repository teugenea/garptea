package handler

import (
	"garptea/ws"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

func WsHandler(c *websocket.Conn) {

	claims := c.Locals("claims").(jwt.MapClaims)
	connectionParams := ws.ConnectionParams{
		UserId:     claims["id"].(string),
		Connection: c,
	}

	defer ws.Unregister(connectionParams)

	ws.Register(connectionParams)

	c.SetPongHandler(func(appData string) error {
		log.Info("pong")
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
