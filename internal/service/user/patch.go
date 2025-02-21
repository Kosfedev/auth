package user

import (
	"context"

	"github.com/Kosfedev/auth/internal/model"
)

func (s *serv) Patch(ctx context.Context, userData *model.UpdatedUserData, id int64) (*model.UserData, error) {
	return s.userRepository.Patch(ctx, userData, id)
}
