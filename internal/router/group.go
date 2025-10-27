package router

import (
	api "github.com/1111mp/gin-app/internal/api/v1"
	"github.com/gin-gonic/gin"
)

// RouterGroupInter -.
type RouterGroupInter interface {
	RegisterPublicRoutes(group *gin.RouterGroup)
	RegisterPrivateRoutes(group *gin.RouterGroup)
}

// RouterGroup -.
type RouterGroup struct {
	UserRouter UserRouterInter
}

// NewRouterGroup -.
func NewRouterGroup(a *api.ApiGroup) *RouterGroup {

	return &RouterGroup{
		&UserRouter{
			userApi: a.UserApi,
		},
	}
}

// RegisterPublicRoutes -.
func (r *RouterGroup) RegisterPublicRoutes(group *gin.RouterGroup) {
	// user
	{
		r.UserRouter.RegisterPublicRoutes(group)
	}
}

// RegisterPublicRoutes -.
func (r *RouterGroup) RegisterPrivateRoutes(group *gin.RouterGroup) {
	// user
	{
		r.UserRouter.RegisterPrivateRoutes(group)
	}
}
