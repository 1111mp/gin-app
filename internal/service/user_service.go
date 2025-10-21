package service

import "github.com/1111mp/gin-app/pkg/logger"

// UserService -.
type UserService struct {
	l logger.Interface
}

func (u *UserService) CreateUser() {
	u.l.Info("CreateUser called")
}
