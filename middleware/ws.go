package middleware

import (
	"fmt"
	"garptea/auth"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func WsUpgrader(c *fiber.Ctx) error {
	queries := c.Queries()
	rawToken := queries["t"]

	if len(rawToken) == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token, err := auth.ParseJwtToken(rawToken)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	claims := token.Claims.(jwt.MapClaims)
	groups := convertToStringArr(claims["groups"].([]interface{}))
	group := groups[auth.ROLE_USER]
	if len(group) == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return c.SendStatus(fiber.StatusUpgradeRequired)
}

func convertToStringArr(data []interface{}) map[string]string {
	s := make(map[string]string, len(data))
	for _, v := range data {
		value := fmt.Sprint(v)
		s[value] = fmt.Sprint(value)
	}
	return s
}
