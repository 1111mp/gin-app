package openapi_service

import (
	"github.com/1111mp/gin-app/internal/repository"
	"github.com/1111mp/gin-app/pkg/jwt"
	"github.com/1111mp/gin-app/pkg/logger"
)

// ServiceGroup -.
type ServiceGroup struct {
}

// NewServiceGroup -.
func NewServiceGroup(r *repository.RepositoryGroup, j jwt.JWTManagerInterface, l logger.Interface) *ServiceGroup {
	return &ServiceGroup{}
}
