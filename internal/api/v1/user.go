package v1

import (
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewUserRoutes(group *gin.RouterGroup, l logger.Interface) {
	userGroup := group.Group("/user")
	{
		userGroup.GET("/:id", GetById)
	}
}

func GetById(c *gin.Context) {

}
