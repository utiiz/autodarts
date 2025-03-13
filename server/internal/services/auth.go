package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/utiiz/autodarts/internal/models"
	"github.com/utiiz/autodarts/internal/repositories"
)

func AuthenticateUser(email, password string) (models.User, error) {
	// Get user by email
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return models.User{}, errors.New("invalid email or password")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, errors.New("invalid email or password")
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
