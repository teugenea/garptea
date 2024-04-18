package middleware

import (
	"garptea/auth"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func WsUpgrader(c *fiber.Ctx) error {
	queries := c.Queries()
	rawToken := queries["t"]

	if len(rawToken) == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims, err := auth.ParseJwtToken(rawToken)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	groups := convertToStringArr(claims.Groups)
	group := groups[auth.ROLE_USER]
	if len(group) == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("claims", claims)
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return c.SendStatus(fiber.StatusUpgradeRequired)
}

func convertToStringArr(data []string) map[string]string {
	s := make(map[string]string, len(data))
	for _, v := range data {
		s[v] = v
	}
	return s
}
