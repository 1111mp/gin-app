package service

import (
	"github.com/1111mp/gin-app/internal/repository"
	"github.com/1111mp/gin-app/pkg/logger"
)

// ServiceGroup -.
type ServiceGroup struct {
	UserService *UserService
}

// NewServiceGroup -.
func NewServiceGroup(r *repository.RepositoryGroup, l logger.Interface) *ServiceGroup {
	return &ServiceGroup{
		UserService: &UserService{
			l:   l,
			rep: r.UserRepository,
		},
	}
}
