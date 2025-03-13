package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/utiiz/autodarts/internal/models"
	"github.com/utiiz/autodarts/internal/repositories"
)

func (h *Handler) SettingsPage(c fiber.Ctx) error {
	return h.RenderComponent(c, "pages/settings.html", nil, fiber.Map{
		"Title": "Settings",
	})
}

func (h *Handler) UpdateSettings(c fiber.Ctx) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

	// Get user from context
	user := c.Locals("user").(models.User)

	// Update user
	user.Name = name
	user.Email = email

	// Save user
	err := repositories.UpdateUser(user)
	if err != nil {
		return h.RenderComponent(c, "pages/settings.html", err, fiber.Map{
			"Title": "Settings",
			"User":  user,
		})
	}

	return h.RenderComponent(c, "pages/settings.html", nil, fiber.Map{
		"Title":   "Settings",
		"User":    user,
		"Success": "Settings updated successfully",
	})
}
