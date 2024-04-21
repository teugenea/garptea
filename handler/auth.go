package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"garptea/auth"
	"garptea/config"
)

func GetJwtToken(c *fiber.Ctx) error {
	queries := c.Queries()
	token, err := auth.GetTokenByAccessCode(queries["code"], queries["state"])
	if err != nil {
		return err
	}
	return c.SendString(token.AccessToken)
}

type hookBody struct {
	ExtendedUser extendedUser `json:"extendedUser"`
}

type extendedUser struct {
	Id string `json:"id"`
}

func ActualizeUserAuth(c *fiber.Ctx) error {
	authHeader := c.Get(auth.HEADER_AUTH)
	if len(authHeader) == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if authHeader == config.GetStringOrEmpty(config.HOOK_SECRET) {
		payload := new(hookBody)
		if err := c.BodyParser(&payload); err != nil {
			log.Errorf("cannot trigger hook: %s", err)
		}
		log.Infof("user was edited id=%s", payload.ExtendedUser.Id)
		auth.ResetUserValidation(payload.ExtendedUser.Id)
	}
	return c.Next()
}
