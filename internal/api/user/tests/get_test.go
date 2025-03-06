package tests

import (
	"context"
	"fmt"
	"testing"

	userImpl "github.com/Kosfedev/auth/internal/api/user"
	modelService "github.com/Kosfedev/auth/internal/model"
	"github.com/Kosfedev/auth/internal/service"
	serviceMocks "github.com/Kosfedev/auth/internal/service/mocks"
	modelHTTP "github.com/Kosfedev/auth/pkg/user_v1/http/types"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		serviceErr = fmt.Errorf("service error")

		req = id

		userData = &modelService.UserData{
			Name:      name,
			Email:     email,
			Role:      0,
			CreatedAt: createdAt,
			UpdatedAt: &updatedAt,
		}

		res = &modelHTTP.ResponseUserData{
			Name:      name,
			Email:     email,
			Role:      0,
			CreatedAt: createdAt,
			UpdatedAt: &updatedAt,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *modelHTTP.ResponseUserData
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(userData, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			userServiceMock := test.userServiceMock(mc)
			api := userImpl.NewImplementation(userServiceMock)

			userData, err := api.Get(test.args.ctx, test.args.req)
			require.Equal(t, test.err, err)
			require.Equal(t, test.want, userData)
		})
	}
}
