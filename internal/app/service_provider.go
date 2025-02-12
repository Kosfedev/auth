package app

import (
	"context"
	"log"

	"github.com/Kosfedev/auth/internal/api/user"
	"github.com/Kosfedev/auth/internal/client/db"
	"github.com/Kosfedev/auth/internal/client/db/pg"
	"github.com/Kosfedev/auth/internal/closer"
	"github.com/Kosfedev/auth/internal/config"
	"github.com/Kosfedev/auth/internal/repository"
	userRepository "github.com/Kosfedev/auth/internal/repository/user"
	userService "github.com/Kosfedev/auth/internal/service"
)

type serviceProvider struct {
	pgConfig config.PGConfig

	dbClient db.Client

	userRepository repository.UserRepository
	userService    userService.UserService
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

		// TODO: add Ping
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) userService.UserService {
	if s.userService == nil {
		s.userService = userService.UserService(
			s.UserRepository(ctx),
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
