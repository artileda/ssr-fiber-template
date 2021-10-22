package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"trovita/configs"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/gofiber/template/mustache"
)

func main() {

	timeoutDefault, _ := time.ParseDuration("1m")
	engine := mustache.NewFileSystem(http.Dir("./views"), ".mustache")

	app := fiber.New(fiber.Config{
		ReadTimeout:  timeoutDefault,
		WriteTimeout: timeoutDefault,
		Views:        engine,
	})

	app.Use(limiter.New(limiter.Config{
		Max:        30,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "request is only allowed 15 each minute"})
		},
	}))

	app.Mount("/", configs.PublicRoutes(app))

	if os.Getenv("ENV") == "development" {
		StartServer(app)
	} else {
		StartServerWithGracefulShutdown(app)
	}
}

// StartServerWithGracefulShutdown function for starting server with a graceful shutdown.
func StartServerWithGracefulShutdown(a *fiber.App) {
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := a.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	// Run server.
	if err := a.Listen(os.Getenv("HOST") + ":" + os.Getenv("PORT")); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

// StartServer func for starting a simple server.
func StartServer(a *fiber.App) {
	// Run server.
	if err := a.Listen(os.Getenv("HOST") + ":" + os.Getenv("PORT")); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
