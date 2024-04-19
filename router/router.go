package router

import (
	"garptea/handler"
	"garptea/middleware"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Get("/ws", middleware.WsUpgrader, websocket.New(handler.WsHandler))
	app.Post("/hook", handler.ActualizeUserAuth)
}

func restricted(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*casdoorsdk.Claims)
	return c.SendString("Welcome " + claims.User.Name)
}
