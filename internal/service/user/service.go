package user

import (
	"github.com/Kosfedev/auth/internal/client/db"
	"github.com/Kosfedev/auth/internal/repository"
	"github.com/Kosfedev/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

// NewService is...
func NewService(userRepository repository.UserRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
