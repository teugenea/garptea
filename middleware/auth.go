package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	"garptea/auth"
	"garptea/handler"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		JWKSetURLs:   []string{auth.GetJwksUrl()},
		ErrorHandler: handler.JwtErrorHandler,
	})
}
