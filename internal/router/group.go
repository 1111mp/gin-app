package router

import (
	api "github.com/1111mp/gin-app/internal/api/v1"
)

// RouterGroup -.
type RouterGroup struct {
	UserRouter *UserRouter
}

// NewRouterGroup -.
func NewRouterGroup(a *api.ApiGroup) *RouterGroup {
	return &RouterGroup{
		UserRouter: &UserRouter{
			userApi: a.UserApi,
		},
	}
}
