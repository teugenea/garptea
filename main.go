package main

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"garptea/config"
	"garptea/router"
	"garptea/ws"
)

func main() {
	config.LoadConfig()
	app := fiber.New()
	go ws.RunWsLoop()
	router.SetupRoutes(app)

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

func createTlsListener() (net.Listener, error) {
	certFile := config.GetStringOrEmpty(config.TLS_CERT_FILE)
	keyFile := config.GetStringOrEmpty(config.TLS_KEY_FILE)
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
