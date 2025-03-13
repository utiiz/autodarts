package handlers

import (
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) HomePage(c fiber.Ctx) error {
	return h.RenderComponent(c, "pages/home.html", nil, fiber.Map{
		"Title": "Home",
	})
}
