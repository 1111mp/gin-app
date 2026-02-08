package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonController struct{}

func NewCommonController() *CommonController {
	return &CommonController{}
}

func (cc *CommonController) Healthz(c *gin.Context) {
	c.Status(http.StatusOK)
}
