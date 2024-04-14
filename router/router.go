package router

import (
	"garptea/handler"
	"garptea/middleware"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
)

func SetupRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://psyduck.home",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/token", handler.GetJwtToken)
	app.Get("/restricted", middleware.Protected(), restricted)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Use(middleware.WsUpgrader)
	app.Get("/ws", websocket.New(handler.WsHandler))
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
