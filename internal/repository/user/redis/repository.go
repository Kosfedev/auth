package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/Kosfedev/auth/internal/repository"
	redigo "github.com/gomodule/redigo/redis"

	"github.com/Kosfedev/auth/internal/client/cache"
	modelService "github.com/Kosfedev/auth/internal/model"
	"github.com/Kosfedev/auth/internal/repository/user/redis/converter"
	modelRepo "github.com/Kosfedev/auth/internal/repository/user/redis/model"
)

// TODO: нарушение DRY? Хочется мапить с типа UserData
const (
	redisFieldID        = "id"
	redisFieldName      = "name"
	redisFieldEmail     = "email"
	redisFieldRole      = "role"
	redisFieldCreatedAt = "created_at"
	redisFieldUpdatedAt = "updated_at"
	redisExpire         = time.Duration(5) * time.Minute
)

var _ repository.UserCacheRepository = (*repo)(nil)

type repo struct {
	cl cache.RedisClient
}

func NewRepository(cl cache.RedisClient) *repo {
	return &repo{cl: cl}
}

func (r *repo) Create(ctx context.Context, userData *modelService.UserData) (int64, error) {
	repoUserData := converter.UserDataToRepo(userData)

	idStr := strconv.FormatInt(userData.ID, 10)
	err := r.cl.HSet(ctx, idStr, repoUserData)
	if err != nil {
		return 0, err
	}
	err = r.cl.Expire(ctx, idStr, redisExpire)
	if err != nil {
		return 0, err
	}

	return userData.ID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*modelService.UserData, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, modelService.ErrorUserNotFound
	}

	var user modelRepo.UserData
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.UserDataFromRepo(&user), nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	idStr := strconv.FormatInt(id, 10)

	return r.cl.HDel(ctx, idStr, redisFieldID, redisFieldName, redisFieldEmail, redisFieldRole, redisFieldCreatedAt, redisFieldUpdatedAt)
}
