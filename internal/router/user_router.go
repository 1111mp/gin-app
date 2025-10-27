package router

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
	userApi *api.UserApi
}

// RegisterPublicRoutes -.
func (u *UserRouter) RegisterPublicRoutes(group *gin.RouterGroup) {
	// userGroup := group.Group("/user")
	// {
	// 	//
	// }
}

// RegisterPrivateRoutes -.
func (u *UserRouter) RegisterPrivateRoutes(group *gin.RouterGroup) {
	userGroup := group.Group("/user")
	{
		userGroup.GET("/:id", u.userApi.GetById)
	}
}
