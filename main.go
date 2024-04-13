package main

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"

	"garptea/auth"
	"garptea/config"
	"garptea/handler"
)

func main() {
	config.LoadConfig()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	app.Get("/token", handler.GetJwtToken)

	app.Use(jwtware.New(jwtware.Config{
		JWKSetURLs:   []string{auth.GetJwksUrl()},
		ErrorHandler: handler.JwtErrorHandler,
	}))
	app.Get("/restricted", restricted)

	if !config.GetBoolOrDefault(config.TLS_ENABLED, true) {
		port := ":" + config.GetStringOrDefault(config.PORT, "3000")
		log.Fatal(app.Listen(port))
	} else {
		listener, err := createTlsListener()
		if err != nil {
			panic(err)
		}
		log.Fatal(app.Listener(listener))
	}
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}

func createTlsListener() (net.Listener, error) {
	certFile := config.GetStringOrDefault(config.TLS_CERT_FILE, "")
	keyFile := config.GetStringOrDefault(config.TLS_KEY_FILE, "")
	if len(certFile) == 0 {
		return nil, errors.New("cannot load tls certificate file")
	}
	if len(keyFile) == 0 {
		return nil, errors.New("cannot load tls key file")
	}
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
	port := ":" + config.GetStringOrDefault(config.PORT, "3000")
	ln, err := tls.Listen("tcp", port, tlsConfig)
	if err != nil {
		return nil, err
	}
	return ln, nil
}
