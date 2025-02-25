package tests

import (
	"context"
	"fmt"
	"testing"

	modelService "github.com/Kosfedev/auth/internal/model"
	"github.com/Kosfedev/auth/internal/repository"
	repositoryMocks "github.com/Kosfedev/auth/internal/repository/mocks"
	service "github.com/Kosfedev/auth/internal/service/user"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userRepoMockFunc func(mc *minimock.Controller) repository.UserRepository
	type args struct {
		ctx context.Context
		req *modelService.NewUserData
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, true, 7)

		serviceErr = fmt.Errorf("service error")

		req = &modelService.NewUserData{
			Name:            name,
			Email:           email,
			Role:            0,
			Password:        password,
			PasswordConfirm: password,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		want         int64
		err          error
		userRepoMock userRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  serviceErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			userRepoMock := test.userRepoMock(mc)
			// TODO: add transactorMock
			serviceMock := service.NewService(userRepoMock)

			newID, err := serviceMock.Create(test.args.ctx, test.args.req)
			require.Equal(t, test.err, err)
			require.Equal(t, test.want, newID)
		})
	}
}
