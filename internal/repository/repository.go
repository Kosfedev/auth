package repository

import (
	"context"

	"github.com/Kosfedev/auth/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, userData *model.NewUserData) (int64, error)
	Get(ctx context.Context, id int64) (*model.UserData, error)
	Patch(ctx context.Context, userData *model.UpdatedUserData, id int64) (*model.UserData, error)
	Delete(ctx context.Context, id int64) error
}
