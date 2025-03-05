package user

import (
	"context"

	"github.com/Kosfedev/auth/internal/model"
)

// TODO: добавить трансактор?
func (s *serv) Patch(ctx context.Context, userData *model.UpdatedUserData, id int64) (*model.UserData, error) {
	updatedData, err := s.userRepository.Patch(ctx, userData, id)
	if err != nil {
		return nil, err
	}

	_, err = s.userRepositoryCache.Create(ctx, updatedData)

	return updatedData, err
}
