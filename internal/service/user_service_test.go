package service_test

import (
	"testing"

	"github.com/1111mp/gin-app/internal/service"
	"github.com/1111mp/gin-app/pkg/logger"
	"go.uber.org/mock/gomock"
)

type test struct {
	name string
	mock func()
	res  interface{}
	err  error
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockUserRepositoryInter(ctrl)
	l := logger.New("", "debug")
	userService := service.NewUserService(l, repo)

	tests := []test{
		{
			name: "create one",
			mock: func() {
				repo.EXPECT().CreateOne().Times(1)
			},
			res: nil,
			err: nil,
		},
	}

	for _, tc := range tests { //nolint:paralleltest // data races here
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			userService.CreateOne()
		})
	}
}
