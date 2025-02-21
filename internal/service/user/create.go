package user

import (
	"context"

	"github.com/Kosfedev/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, userData *model.NewUserData) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.Create(ctx, userData)
		if errTx != nil {
			return errTx
		}

		// дополнительный запрос в базу для демонстрации трансакции
		// TODO: убрать, когда появится реальная логика для трансакции
		_, errTx = s.userRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
