package app

import (
	"context"
	"log"

	"github.com/Kosfedev/auth/internal/api/user"
	"github.com/Kosfedev/auth/internal/client/db"
	"github.com/Kosfedev/auth/internal/client/db/pg"
	"github.com/Kosfedev/auth/internal/client/db/transaction"
	"github.com/Kosfedev/auth/internal/closer"
	"github.com/Kosfedev/auth/internal/config"
	"github.com/Kosfedev/auth/internal/repository"
	userRepository "github.com/Kosfedev/auth/internal/repository/user"
	userServInterface "github.com/Kosfedev/auth/internal/service"
	userService "github.com/Kosfedev/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig config.PGConfig

	dbClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository
	userService    userServInterface.UserService
	userImpl       *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		return cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) userServInterface.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}
