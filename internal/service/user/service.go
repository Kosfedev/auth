package user

import (
	"github.com/Kosfedev/auth/internal/client/db"
	"github.com/Kosfedev/auth/internal/repository"
	"github.com/Kosfedev/auth/internal/service"
)

type serv struct {
	userRepository      repository.UserRepository
	userRepositoryCache repository.UserCacheRepository
	txManager           db.TxManager
}

// NewService is...
func NewService(userRepository repository.UserRepository, userRepositoryCache repository.UserCacheRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepository:      userRepository,
		userRepositoryCache: userRepositoryCache,
		txManager:           txManager,
	}
}
