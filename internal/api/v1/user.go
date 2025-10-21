package v1

import (
	"github.com/1111mp/gin-app/internal/service"
	"github.com/gin-gonic/gin"
)

// UserApi -.
type UserApi struct {
	userService *service.UserService
}

func (u *UserApi) GetById(c *gin.Context) {
	u.userService.CreateUser()
}
