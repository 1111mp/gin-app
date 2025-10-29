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
	PostRouter PostRouterInter
}

// NewRouterGroup -.
func NewRouterGroup(a *api.ApiGroup) *RouterGroup {

	return &RouterGroup{
		&UserRouter{
			userApi: a.UserApi,
		},
		&PostRouter{
			postApi: a.PostApi,
		},
	}
}

// RegisterPublicRoutes -.
func (r *RouterGroup) RegisterPublicRoutes(group *gin.RouterGroup) {
	// users
	{
		r.UserRouter.RegisterPublicRoutes(group)
	}
	// posts
	{
		r.PostRouter.RegisterPublicRoutes(group)
	}
}

// RegisterPublicRoutes -.
func (r *RouterGroup) RegisterPrivateRoutes(group *gin.RouterGroup) {
	// users
	{
		r.UserRouter.RegisterPrivateRoutes(group)
	}
	// posts
	{
		r.PostRouter.RegisterPrivateRoutes(group)
	}
}
