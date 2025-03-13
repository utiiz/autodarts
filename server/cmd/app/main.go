package main

import (
	"embed"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/utiiz/autodarts/internal/handlers"
	"github.com/utiiz/autodarts/internal/middlewares"
)

//go:embed ../../templates/layouts/*.html ../../templates/components/*.html ../../templates/pages/*.html
var templates embed.FS

var t = template.Must(template.ParseFS(templates, "templates/layouts/*.html", "templates/components/*.html", "templates/pages/*.html"))

func main() {
	// Parse templates
	templates, err := handlers.ParseFiles()
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// Create handler
	h := &handlers.Handler{
		Templates: templates,
	}

	// Create fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// // Static files
	// app.Static("/static", "./static")

	// Auth routes
	app.Get("/login", h.LoginPage)
	app.Post("/login", h.Login)
	app.Get("/logout", h.Logout)

	// Protected routes
	app.Use(middlewares.AuthMiddleware())
	app.Get("/", h.HomePage)
	app.Get("/settings", h.SettingsPage)
	app.Post("/settings", h.UpdateSettings)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
