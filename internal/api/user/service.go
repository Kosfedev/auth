package user

import (
	"github.com/Kosfedev/auth/internal/service"
)

type Implementation struct {
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
