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

func TestCreate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *modelHTTP.RequestNewUserData
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, true, 7)

		serviceErr = fmt.Errorf("service error")

		req = &modelHTTP.RequestNewUserData{
			Name:            name,
			Email:           email,
			Role:            0,
			Password:        password,
			PasswordConfirm: password,
		}

		newUserData = &modelService.NewUserData{
			Name:            name,
			Email:           email,
			Role:            0,
			Password:        password,
			PasswordConfirm: password,
		}

		res = &modelHTTP.ResponseUserID{
			ID: id,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *modelHTTP.ResponseUserID
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
				mock.CreateMock.Expect(ctx, newUserData).Return(id, nil)
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
				mock.CreateMock.Expect(ctx, newUserData).Return(0, serviceErr)
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

			newID, err := api.Create(test.args.ctx, test.args.req)
			require.Equal(t, test.err, err)
			require.Equal(t, test.want, newID)
		})
	}
}
