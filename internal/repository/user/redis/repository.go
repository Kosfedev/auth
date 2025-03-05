package redis

import (
	"context"
	"strconv"

	"github.com/Kosfedev/auth/internal/repository"
	redigo "github.com/gomodule/redigo/redis"

	"github.com/Kosfedev/auth/internal/client/cache"
	modelService "github.com/Kosfedev/auth/internal/model"
	"github.com/Kosfedev/auth/internal/repository/user/redis/converter"
	modelRepo "github.com/Kosfedev/auth/internal/repository/user/redis/model"
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

	return userData.ID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*modelService.UserData, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, modelService.ErrorNoteNotFound
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

	// TODO: нарушение DRY
	err := r.cl.HDel(ctx, idStr, "id", "name", "email", "role", "created_at", "updated_at")

	return err
}

// TODO:убрать методы заглушки (не нужны для кеширования)
func (r *repo) Patch(ctx context.Context, userData *modelService.UpdatedUserData, id int64) (*modelService.UserData, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, modelService.ErrorNoteNotFound
	}

	var user modelRepo.UserData
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.UserDataFromRepo(&user), nil
}
