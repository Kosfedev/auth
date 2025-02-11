package user

import (
	"github.com/Kosfedev/auth/internal/repository"
	"github.com/Kosfedev/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) service.UserService {
	return &serv{
		userRepository: userRepository,
	}
}
