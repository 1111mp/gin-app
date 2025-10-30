package api_router

import (
	api "github.com/1111mp/gin-app/internal/api/v1"
	"github.com/gin-gonic/gin"
)

// RouterGroupInter -.
type ApiRouterGroupInter interface {
	RegisterPublicRoutes(group *gin.RouterGroup)
	RegisterPrivateRoutes(group *gin.RouterGroup)
}

// RouterGroup -.
type ApiRouterGroup struct {
	UserRouter        UserRouterInter
	PostRouter        PostRouterInter
	AccessTokenRouter AccessTokenRouterInter
}

// NewRouterGroup -.
func NewRouterGroup(a *api.ApiGroup) *ApiRouterGroup {

	return &ApiRouterGroup{
		&UserRouter{
			userApi: a.UserApi,
		},
		&PostRouter{
			postApi: a.PostApi,
		},
		&AccessTokenRouter{
			accessTokenApi: a.AccessTokenApi,
		},
	}
}

// RegisterPublicRoutes -.
func (r *ApiRouterGroup) RegisterPublicRoutes(group *gin.RouterGroup) {
	// users
	{
		r.UserRouter.RegisterPublicRoutes(group)
	}
	// posts
	{
		r.PostRouter.RegisterPublicRoutes(group)
	}
	// access-tokens
	{
		r.AccessTokenRouter.RegisterPublicRoutes(group)
	}
}

// RegisterPublicRoutes -.
func (r *ApiRouterGroup) RegisterPrivateRoutes(group *gin.RouterGroup) {
	// users
	{
		r.UserRouter.RegisterPrivateRoutes(group)
	}
	// posts
	{
		r.PostRouter.RegisterPrivateRoutes(group)
	}
	// access-tokens
	{
		r.AccessTokenRouter.RegisterPrivateRoutes(group)
	}
}
