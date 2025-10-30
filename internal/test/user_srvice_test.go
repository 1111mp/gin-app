package test

import (
	"context"
	"testing"

	ent "github.com/1111mp/gin-app/ent"
	dto "github.com/1111mp/gin-app/internal/dto"
	api_service "github.com/1111mp/gin-app/internal/service/api"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type test struct {
	name  string
	mock  func()
	res   interface{}
	err   error
	token string
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockUserRepositoryInter(ctrl)
	jwt := NewMockJWTManagerInterface(ctrl)
	l := logger.New("", "debug")
	userService := api_service.NewUserService(l, repo, jwt)

	ctx := context.Background()
	inputDto := dto.UserCreateOneDto{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "verysecurepwd12",
	}
	mockUser := &ent.User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}

	tests := []test{
		{
			name: "create one empty result",
			mock: func() {
				repo.EXPECT().CreateOne(ctx, inputDto).Return(mockUser, nil)
				jwt.EXPECT().GenerateToken(mockUser.ID).Return("mock-token", nil)
			},
			res:   mockUser,
			token: "mock-token",
			err:   nil,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			res, token, err := userService.CreateOne(ctx, inputDto)

			require.Equal(t, res.ID, mockUser.ID)
			require.Equal(t, token, localTc.token)
			require.ErrorIs(t, err, localTc.err)
		})
	}
}
