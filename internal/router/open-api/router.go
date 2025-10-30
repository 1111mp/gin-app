package openapi_router

import (
	openapi_v1 "github.com/1111mp/gin-app/internal/open-api/v1"
	"github.com/gin-gonic/gin"
)

// ApiRouterGroupInter -.
type ApiRouterGroupInter interface {
	RegisterRoutes(group *gin.RouterGroup)
}

// ApiRouterGroup -.
type ApiRouterGroup struct {
}

// NewRouterGroup -.
func NewRouterGroup(a *openapi_v1.ApiGroup) *ApiRouterGroup {
	return &ApiRouterGroup{}
}

// RegisterRoutes -.
func (o *ApiRouterGroup) RegisterRoutes(group *gin.RouterGroup) {

}
