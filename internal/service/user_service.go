package service

import (
	"github.com/1111mp/gin-app/internal/repository"
	"github.com/1111mp/gin-app/pkg/logger"
)

// UserService -.
type UserService struct {
	l   logger.Interface
	rep *repository.UserRepository
}

// CreateUser -.
func (u *UserService) CreateUser() {
	u.l.Info("CreateUser called")
	u.rep.CreateOne()
}
