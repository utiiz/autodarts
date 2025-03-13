package main

import (
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	_ "github.com/lib/pq"
)

func main() {
	app := fiber.New()

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	log.Println("GET /")
	// 	return nil
	// })

	listenAddr := ":8080"
	slog.Info("HTTP Server started", "listenAddr", listenAddr)
	log.Fatal(app.Listen(listenAddr))
}
