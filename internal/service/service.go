package service

import (
	"github.com/1111mp/gin-app/internal/repository"
	api_service "github.com/1111mp/gin-app/internal/service/api"
	openapi_service "github.com/1111mp/gin-app/internal/service/open-api"
	"github.com/1111mp/gin-app/pkg/jwt"
	"github.com/1111mp/gin-app/pkg/logger"
)

// NewServiceGroup -.
func NewServiceGroup(
	r *repository.RepositoryGroup,
	j jwt.JWTManagerInterface,
	l logger.Interface,
) (*api_service.ServiceGroup, *openapi_service.ServiceGroup) {
	return api_service.NewServiceGroup(r, j, l), openapi_service.NewServiceGroup(r, j, l)
}
