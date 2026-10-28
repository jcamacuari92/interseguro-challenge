package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/interseguro/go-api/internal/clients"
	"github.com/interseguro/go-api/internal/config"
	"github.com/interseguro/go-api/internal/handlers"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		AppName:      "Interseguro Go API v1.0",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	h := handlers.New(clients.NewNodeClient(cfg.NodeAPIURL))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "go-api"})
	})

	api := app.Group("/api/v1")
	api.Post("/qr", h.QR)
	api.Post("/rotate", h.Rotate)

	go func() {
		log.Printf("🚀 Go API escuchando en :%s (Node API: %s)", cfg.Port, cfg.NodeAPIURL)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("error al iniciar servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Apagando Go API...")
	_ = app.ShutdownWithTimeout(5 * time.Second)
}