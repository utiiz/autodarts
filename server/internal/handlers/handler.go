package handlers

import (
	"html/template"
	"log"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	Templates *template.Template
}

// ParseFiles loads and parses all HTML templates
func ParseFiles() (*template.Template, error) {
	// Parse all template files
	tmpl := template.New("")

	// Define template functions
	tmpl = tmpl.Funcs(template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	})

	// Parse templates from layouts, components, and pages
	layouts, err := filepath.Glob("templates/layouts/*.html")
	if err != nil {
		return nil, err
	}

	components, err := filepath.Glob("templates/components/*.html")
	if err != nil {
		return nil, err
	}

	pages, err := filepath.Glob("templates/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Combine all templates
	templateFiles := append(layouts, components...)
	templateFiles = append(templateFiles, pages...)

	tmpl, err = tmpl.ParseFiles(templateFiles...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// RenderComponent renders a template with the provided data
func (h *Handler) RenderComponent(c fiber.Ctx, templateName string, err error, data fiber.Map) error {
	// If there's an error, add it to the data
	if err != nil {
		log.Printf("Error: %v", err)
		if data == nil {
			data = fiber.Map{}
		}
		data["Error"] = err.Error()
	}

	// Add default data
	if data == nil {
		data = fiber.Map{}
	}

	// Check if user is authenticated
	if auth := c.Locals("authenticated"); auth != nil {
		data["Authenticated"] = auth
		data["User"] = c.Locals("user")
	}

	// Determine if this is an HTMX request
	isHtmx := c.Get("HX-Request") == "true"

	// If it's an HTMX request, render just the component
	if isHtmx && !strings.HasPrefix(templateName, "layouts/") {
		return h.Templates.ExecuteTemplate(c.Response().BodyWriter(), templateName, data)
	}

	// Otherwise, render with the base layout
	data["Content"] = templateName
	log.Printf("Data: %v", data)
	log.Printf("Content: %v", data["Content"])
	return h.Templates.ExecuteTemplate(c.Response().BodyWriter(), "layouts/base.html", data)
}
