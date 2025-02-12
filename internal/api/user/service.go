package user

import (
	"github.com/Kosfedev/auth/internal/service"
)

// Implementation is...
type Implementation struct {
	userService service.UserService
}

// NewImplementation is...
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
