package handler

import (
	"github.com/gofiber/fiber/v2"

	"garptea/auth"
)

func JwtErrorHandler(c *fiber.Ctx, err error) error {
	return c.Redirect(auth.GetLoginUrl())
}

func GetJwtToken(c *fiber.Ctx) error {
	queries := c.Queries()
	token := auth.GetTokenByAccessCode(queries["code"], queries["state"])
	c.Locals("user", token)
	return c.SendString(token)
}
