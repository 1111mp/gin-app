package router

import (
	api "github.com/1111mp/gin-app/internal/api/v1"
	"github.com/1111mp/gin-app/pkg/utils"
	"github.com/gin-gonic/gin"
)

// PostRouterInter -.
type PostRouterInter interface {
	RegisterPublicRoutes(group *gin.RouterGroup)
	RegisterPrivateRoutes(group *gin.RouterGroup)
}

// PostRouter -.
type PostRouter struct {
	postApi api.PostApiInter
}

// RegisterPublicRoutes -.
func (p *PostRouter) RegisterPublicRoutes(group *gin.RouterGroup) {
	// postGroup := group.Group("/posts")
	// {

	// }
}

// RegisterPrivateRoutes -.
func (p *PostRouter) RegisterPrivateRoutes(group *gin.RouterGroup) {
	postGroup := group.Group("/posts")
	{
		postGroup.POST("", utils.HandlerWithUser(p.postApi.CreateOne))
		postGroup.GET("/:id", p.postApi.GetById)
	}
}
