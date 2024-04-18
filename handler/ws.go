package handler

import (
	"garptea/ws"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

func WsHandler(c *websocket.Conn) {
	defer ws.Unregister(c)

	claims := c.Locals("claims").(jwt.MapClaims)
	ws.Register(ws.ConnectionParams{
		UserId:     claims["id"].(string),
		Connection: c,
	})
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("read error: %s", err)
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
