package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	"garptea/handler"
	"garptea/util"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		JWKSetURLs:   []string{util.GetJwksUrl()},
		ErrorHandler: handler.JwtErrorHandler,
	})
}
