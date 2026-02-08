package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CommonController struct{}

func NewCommonController() *CommonController {
	return &CommonController{}
}

func (cc *CommonController) Healthz(c *gin.Context) {
	time.Sleep(5 * time.Second)
	c.Status(http.StatusOK)
}
