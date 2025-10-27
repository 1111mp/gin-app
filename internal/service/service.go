package service

import (
	"github.com/1111mp/gin-app/internal/repository"
	"github.com/1111mp/gin-app/pkg/jwt"
	"github.com/1111mp/gin-app/pkg/logger"
)

// ServiceGroup -.
type ServiceGroup struct {
	UserService UserServiceInter
}

// NewServiceGroup -.
func NewServiceGroup(r *repository.RepositoryGroup, j jwt.JWTManagerInterface, l logger.Interface) *ServiceGroup {
	return &ServiceGroup{
		&UserService{
			l:   l,
			rep: r.UserRepository,
			jwt: j,
		},
	}
}
