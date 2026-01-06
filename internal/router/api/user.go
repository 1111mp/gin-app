package api_router

import (
	api "github.com/1111mp/gin-app/internal/api/v1"
	"github.com/gin-gonic/gin"
)

// UserRouterInter -.
type UserRouterInter interface {
	RegisterPublicRoutes(group *gin.RouterGroup)
	RegisterPrivateRoutes(group *gin.RouterGroup)
}

// UserRouter -.
type UserRouter struct {
	userApi api.UserApiInter
}

// RegisterPublicRoutes -.
func (u *UserRouter) RegisterPublicRoutes(group *gin.RouterGroup) {
	userGroup := group.Group("/users")
	{
		userGroup.POST("", u.userApi.CreateOne)
	}
}

// RegisterPrivateRoutes -.
func (u *UserRouter) RegisterPrivateRoutes(group *gin.RouterGroup) {
	userGroup := group.Group("/users")
	{
		userGroup.GET("/:id", u.userApi.GetById)
	}
}
