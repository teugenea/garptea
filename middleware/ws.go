package middleware

import (
	"garptea/auth"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func WsUpgrader(c *fiber.Ctx) error {
	queries := c.Queries()
	rawToken := queries["t"]

	if len(rawToken) == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims, err := auth.ParseJwtToken(rawToken)
	if err != nil {
		log.Errorf("cannot parse token: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	groups := convertToStringArr(claims.Groups)
	group := groups[auth.ROLE_USER]
	if len(group) == 0 {
		log.Errorf("user %s is not in group %s", claims.User.Id, auth.ROLE_USER)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals(auth.LOCALS_CLAIMS, claims)
	c.Locals(auth.LOCALS_TOKEN, rawToken)
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
