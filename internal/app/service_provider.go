package app

import (
	"context"
	"log"

	"github.com/Kosfedev/auth/internal/api/user"
	"github.com/Kosfedev/auth/internal/closer"
	"github.com/Kosfedev/auth/internal/config"
	"github.com/Kosfedev/auth/internal/repository"
	userRepository "github.com/Kosfedev/auth/internal/repository/user"
	userService "github.com/Kosfedev/auth/internal/service"
	"github.com/jackc/pgx/v5"
)

type serviceProvider struct {
	pgConfig       config.PGConfig
	pgCon          *pgx.Conn
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

func (s *serviceProvider) PGCon(ctx context.Context) *pgx.Conn {
	if s.pgCon == nil {
		dsn := s.PGConfig().DSN()
		con, err := pgx.Connect(ctx, dsn)
		if err != nil {
			log.Fatalf("failed to establish connection to \"%v\": %v", dsn, err.Error())
		}
		closer.Add(func() error {
			return con.Close(ctx)
		})

		s.pgCon = con
	}

	return s.pgCon
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.PGCon(ctx))
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
