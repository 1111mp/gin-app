package v1

import (
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewRoutes(group *gin.RouterGroup, l logger.Interface) {
	NewUserRoutes(group, l)
}
