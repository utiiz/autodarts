package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v4"

	"github.com/utiiz/autodarts/internal/repositories"
)

func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get the JWT token from the cookie
		cookie := c.Cookies("jwt")
		if cookie == "" {
			return c.Redirect().To("/login")
		}

		// Parse the JWT token
		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil // Use a proper secret key from config
		})

		if err != nil || !token.Valid {
			return c.Redirect().To("/login")
		}

		// Get the user ID from the token
		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))

		// Get the user from the database
		user, err := repositories.GetUserByID(userID)
		if err != nil {
			return c.Redirect().To("/login")
		}

		// Set the user in the context
		c.Locals("authenticated", true)
		c.Locals("user", user)

		return c.Next()
	}
}
