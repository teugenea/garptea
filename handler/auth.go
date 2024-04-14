package handler

import (
	"github.com/gofiber/fiber/v2"

	"garptea/util"
)

func JwtErrorHandler(c *fiber.Ctx, err error) error {
	return c.Redirect(util.GetLoginUrl())
}

func GetJwtToken(c *fiber.Ctx) error {
	queries := c.Queries()
	token := util.GetTokenByAccessCode(queries["code"])
	c.Locals("user", token)
	return c.SendString(token)
}
