package user

import (
	"context"
)

// TODO: добавить трансактор?
func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.userCacheRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return s.userRepository.Delete(ctx, id)
}
