package service

import (
	"context"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/internal/repository"
	"github.com/1111mp/gin-app/pkg/logger"
)

// UserService -.
type UserService struct {
	l   logger.Interface
	rep repository.UserRepository
}

// NewUserService -.
func NewUserService(l logger.Interface, rep repository.UserRepository) *UserService {
	return &UserService{l: l, rep: rep}
}

// CreateUser -.
func (u *UserService) CreateUser() {
	u.l.Info("CreateUser called")
	u.rep.CreateOne()
}

// GetById -.
func (u *UserService) GetById(ctx context.Context, id int) (*ent.UserEntity, error) {
	user, err := u.rep.GetById(ctx, id)
	if err != nil {

		return nil, err
	}

	return user.IntoEntity(), nil
}
