package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"garptea/auth"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := extractAuthHeader(c)
		if err != nil {
			return c.Redirect(auth.GetLoginUrl())
		}
		claims, err := auth.ParseJwtToken(token)
		if err != nil {
			return c.Redirect(auth.GetLoginUrl())
		}
		c.Locals(auth.LOCALS_CLAIMS, claims)
		return c.Next()
	}
}

func extractAuthHeader(c *fiber.Ctx) (string, error) {
	header := c.Get("Authorization")
	authScheme := "Bearer"
	l := len(authScheme)
	if len(header) > l+1 && strings.EqualFold(header[:l], authScheme) {
		return strings.TrimSpace(header[l:]), nil
	}
	return "", errors.New("missing or malformed JWT")
}
