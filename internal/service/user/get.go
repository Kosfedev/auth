package user

import (
	"context"

	"github.com/Kosfedev/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.UserData, error) {
	return s.userRepository.Get(ctx, id)
}
