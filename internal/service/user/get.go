package user

import (
	"context"
	"errors"

	"github.com/Kosfedev/auth/internal/model"
)

// TODO: добавить трансактор?
func (s *serv) Get(ctx context.Context, id int64) (*model.UserData, error) {
	userData, err := s.userCacheRepository.Get(ctx, id)
	if errors.Is(err, model.ErrorUserNotFound) {
		userData, err = s.userRepository.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		_, err = s.userCacheRepository.Create(ctx, userData)
		if err != nil {
			return nil, err
		}
	}

	return userData, err
}
