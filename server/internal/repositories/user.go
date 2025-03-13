package repositories

import (
	"errors"
	"sync"

	"github.com/utiiz/autodarts/internal/models"
)

// This is a simple in-memory repository
// In a real application, you would use a database
var (
	users             = make(map[uint]models.User)
	usersByEmail      = make(map[string]models.User)
	userID       uint = 0
	userMutex    sync.RWMutex
)

// Initialize with a test user
func init() {
	// Password would be hashed in a real application
	CreateUser(models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "$2a$10$JZRlmODyikZlZEvUKZE6b.zWdYiMbk7tG9JLEzWQ.0MqE7XkWwouS", // "password"
	})
}

func GetUserByID(id uint) (models.User, error) {
	userMutex.RLock()
	defer userMutex.RUnlock()

	user, exists := users[id]
	if !exists {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	userMutex.RLock()
	defer userMutex.RUnlock()

	user, exists := usersByEmail[email]
	if !exists {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func CreateUser(user models.User) (models.User, error) {
	userMutex.Lock()
	defer userMutex.Unlock()

	// Check if email already exists
	if _, exists := usersByEmail[user.Email]; exists {
		return models.User{}, errors.New("email already in use")
	}

	// Assign ID
	userID++
	user.ID = userID

	// Store user
	users[user.ID] = user
	usersByEmail[user.Email] = user

	return user, nil
}

func UpdateUser(user models.User) error {
	userMutex.Lock()
	defer userMutex.Unlock()

	// Check if user exists
	if _, exists := users[user.ID]; !exists {
		return errors.New("user not found")
	}

	// Update user
	users[user.ID] = user
	usersByEmail[user.Email] = user

	return nil
}
