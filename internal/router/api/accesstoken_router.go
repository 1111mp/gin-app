package api_router

import (
	api "github.com/1111mp/gin-app/internal/api/v1"
	"github.com/1111mp/gin-app/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AccessTokenRouterInter -.
type AccessTokenRouterInter interface {
	RegisterPublicRoutes(group *gin.RouterGroup)
	RegisterPrivateRoutes(group *gin.RouterGroup)
}

// AccessTokenRouter -.
type AccessTokenRouter struct {
	accessTokenApi api.AccessTokenApiInter
}

// RegisterPublicRoutes -.
func (u *AccessTokenRouter) RegisterPublicRoutes(group *gin.RouterGroup) {
	// accessTokenGroup := group.Group("/access-tokens")
	// {
	// }
}

// RegisterPrivateRoutes -.
func (u *AccessTokenRouter) RegisterPrivateRoutes(group *gin.RouterGroup) {
	accessTokenGroup := group.Group("/access-tokens")
	{
		accessTokenGroup.POST("", utils.HandlerWithUser(u.accessTokenApi.CreateOne))
	}
}
