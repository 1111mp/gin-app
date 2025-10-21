package service

import "github.com/1111mp/gin-app/pkg/logger"

// ServiceGroup -.
type ServiceGroup struct {
	UserService *UserService
}

// NewServiceGroup -.
func NewServiceGroup(l logger.Interface) *ServiceGroup {
	return &ServiceGroup{
		UserService: &UserService{
			l: l,
		},
	}
}
