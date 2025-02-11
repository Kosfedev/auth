package user

import (
	"context"

	"github.com/Kosfedev/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, userData *model.NewUserData) (int64, error) {
	return s.userRepository.Create(ctx, userData)
}
