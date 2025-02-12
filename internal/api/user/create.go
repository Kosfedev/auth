package user

import (
	"context"
	"log"

	"github.com/Kosfedev/auth/internal/model"
	modelHTTP "github.com/Kosfedev/auth/pkg/user_v1/http/types"
)

// Create is...
func (i *Implementation) Create(ctx context.Context, userData *model.NewUserData) (*modelHTTP.ResponseUserID, error) {
	id, err := i.userService.Create(ctx, userData)
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %d", id)

	return &modelHTTP.ResponseUserID{
		ID: id,
	}, nil
}
