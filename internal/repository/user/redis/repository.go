package redis

import (
	"context"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/Kosfedev/auth/internal/client/cache"
	modelService "github.com/Kosfedev/auth/internal/model"
	"github.com/Kosfedev/auth/internal/repository/user/redis/converter"
	modelRepo "github.com/Kosfedev/auth/internal/repository/user/redis/model"
)

type repo struct {
	cl cache.RedisClient
}

func NewRepository(cl cache.RedisClient) *repo {
	return &repo{cl: cl}
}

func (r *repo) Create(ctx context.Context, newUserData *modelService.NewUserData) (int64, error) {
	id := int64(1)

	repoUserData := converter.NewUserDataToRepo(newUserData)

	idStr := strconv.FormatInt(id, 10)
	err := r.cl.HSet(ctx, idStr, repoUserData)
	if err != nil {
		return 0, err
	}

	return id, nil
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

// TODO:убрать методы заглушки (не нужны для кеширования)
func (r *repo) Delete(ctx context.Context, id int64) error {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return err
	}

	if len(values) == 0 {
		return modelService.ErrorNoteNotFound
	}

	var user modelRepo.UserData
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return err
	}

	return nil
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
