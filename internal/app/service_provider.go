package app

import (
	"context"
	"log"

	"github.com/Kosfedev/auth/internal/api/user"
	"github.com/Kosfedev/auth/internal/client/cache"
	"github.com/Kosfedev/auth/internal/client/cache/redis"
	"github.com/Kosfedev/auth/internal/client/db"
	"github.com/Kosfedev/auth/internal/client/db/pg"
	"github.com/Kosfedev/auth/internal/client/db/transaction"
	"github.com/Kosfedev/auth/internal/closer"
	"github.com/Kosfedev/auth/internal/config"
	"github.com/Kosfedev/auth/internal/config/env"
	"github.com/Kosfedev/auth/internal/repository"
	userRepositoryPg "github.com/Kosfedev/auth/internal/repository/user/pg"
	userRepositoryRedis "github.com/Kosfedev/auth/internal/repository/user/redis"
	userServInterface "github.com/Kosfedev/auth/internal/service"
	userService "github.com/Kosfedev/auth/internal/service/user"

	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	pgConfig    config.PGConfig
	redisConfig config.RedisConfig

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	userRepository      repository.UserRepository
	userCacheRepository repository.UserCacheRepository
	userService         userServInterface.UserService
	userImpl            *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
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

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		cfg := s.RedisConfig()

		s.redisPool = &redigo.Pool{
			MaxIdle:     cfg.MaxIdle(),
			IdleTimeout: cfg.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", cfg.Address())
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepositoryPg.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserCacheRepository() repository.UserCacheRepository {
	if s.userCacheRepository == nil {
		s.userCacheRepository = userRepositoryRedis.NewRepository(s.RedisClient())
	}

	return s.userCacheRepository
}

func (s *serviceProvider) UserService(ctx context.Context) userServInterface.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.UserCacheRepository(),
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
