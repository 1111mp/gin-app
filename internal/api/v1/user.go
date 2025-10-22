package v1

import (
	"github.com/1111mp/gin-app/internal/service"
	"github.com/gin-gonic/gin"
)

// UserApi -.
type UserApi struct {
	userService *service.UserService
}

// GetById godoc
// @Summary     Get user by ID
// @Description Retrieve user information by given user ID
// @ID          getUserById
// @Tags        user
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Router      /user/{id} [get]
func (u *UserApi) GetById(c *gin.Context) {
	u.userService.CreateUser()
}
