package api_router

import (
	api "github.com/1111mp/gin-app/internal/api/v1"
	"github.com/gin-gonic/gin"
)

// AuthRouterInter -.
type AuthRouterInter interface {
	RegisterPublicRoutes(group *gin.RouterGroup)
	RegisterPrivateRoutes(group *gin.RouterGroup)
}

// AuthRouter -.
type AuthRouter struct {
	authApi api.AuthApiInter
}

// RegisterPublicRoutes -.
func (u *AuthRouter) RegisterPublicRoutes(group *gin.RouterGroup) {
	authGroup := group.Group("/auth")
	{
		authGroup.GET("/login", u.authApi.LoginHandler)
		authGroup.GET("/callback", u.authApi.CallbackHandler)
	}
}

// RegisterPrivateRoutes -.
func (u *AuthRouter) RegisterPrivateRoutes(group *gin.RouterGroup) {

}
