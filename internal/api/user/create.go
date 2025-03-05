package user

import (
	"context"

	"github.com/Kosfedev/auth/internal/converter"
	modelHTTP "github.com/Kosfedev/auth/pkg/user_v1/http/types"
)

// Create is...
func (i *Implementation) Create(ctx context.Context, userData *modelHTTP.RequestNewUserData) (*modelHTTP.ResponseUserID, error) {
	id, err := i.userService.Create(ctx, converter.NewDataFromServiceToHTTP(userData))
	if err != nil {
		return nil, err
	}

	return &modelHTTP.ResponseUserID{
		ID: id,
	}, nil
}
