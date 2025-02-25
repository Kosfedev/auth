package tests

import (
	"context"
	"fmt"
	"testing"

	userImpl "github.com/Kosfedev/auth/internal/api/user"
	"github.com/Kosfedev/auth/internal/service"
	serviceMocks "github.com/Kosfedev/auth/internal/service/mocks"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)
		id  = gofakeit.Int64()

		serviceErr = fmt.Errorf("service error")

		req = id
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
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

			err := api.Delete(test.args.ctx, test.args.req)
			require.Equal(t, test.err, err)
		})
	}
}
