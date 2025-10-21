package router

import (
	api "github.com/1111mp/gin-app/internal/api/v1"
	"github.com/gin-gonic/gin"
)

// UserRouter -.
type UserRouter struct {
	userApi *api.UserApi
}

// RegisterRoutes -.
func (u *UserRouter) RegisterRoutes(group *gin.RouterGroup) {
	userGroup := group.Group("/user")
	{
		userGroup.GET("/:id", u.userApi.GetById)
	}
}
