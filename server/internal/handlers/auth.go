package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/utiiz/autodarts/internal/services"
)

func (h *Handler) LoginPage(c fiber.Ctx) error {
	return h.RenderComponent(c, "pages/login.html", nil, fiber.Map{
		"Title": "Login",
	})
}

func (h *Handler) Login(c fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Your authentication logic here
	// This is a simplified example, replace with your actual authentication service
	user, err := services.AuthenticateUser(email, password)
	if err != nil {
		return h.RenderComponent(c, "pages/login.html", errors.New("Invalid email or password"), fiber.Map{
			"Title": "Login",
		})
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Sign token
	t, err := token.SignedString([]byte("your_secret_key")) // Use a proper secret key from config
	if err != nil {
		return h.RenderComponent(c, "pages/login.html", errors.New("Could not login"), fiber.Map{
			"Title": "Login",
		})
	}

	// Set cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(72 * time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	// Redirect to home
	return c.Redirect().To("/")
}

func (h *Handler) Logout(c fiber.Ctx) error {
	// Clear cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.Redirect().To("/login")
}
