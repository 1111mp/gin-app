package v1

import (
	"net/http"

	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/internal/service"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/response"
	"github.com/gin-gonic/gin"
)

// UserApiInter -.
type UserApiInter interface {
	CreateOne(c *gin.Context)
	GetById(c *gin.Context)
}

// UserApi -.
type UserApi struct {
	cfg         config.ConfigInterface
	userService service.UserServiceInter
}

// CreateOne godoc
// @Summary 		Create a new user
// @Description Creates a new user account
// @ID          CreateOne
// @Tags 				Users
// @Accept 			json
// @Produce 		json
// @Param 			data body dto.CreateOneUserDto true "User data"
// @Success 		200 {object} response.UserAPIResponse "User created successfully"
// @Failure     400 {object} errors.APIError "Bad request (invalid params)"
// @Failure     500 {object} errors.APIError "Internal server error"
// @Router 			/api/v1/users [post]
func (u *UserApi) CreateOne(c *gin.Context) {
	ctx := c.Request.Context()

	var dto dto.CreateOneUserDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(
			errors.NewAPIError(
				http.StatusBadRequest,
				err.Error(),
			),
		)
		return
	}

	user, token, err := u.userService.CreateOne(ctx, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie(u.cfg.HTTP().CookieName, token, 3600, "/", "", true, true)

	response.WriteSuccess(c, user)
}

// GetById godoc
// @Summary     Get user by ID
// @Description Retrieve user information by given user ID
// @ID          getUserById
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Success     200 {object} response.UserAPIResponse "Successfully retrieved user"
// @Failure     400 {object} errors.APIError "Bad request (invalid ID)"
// @Failure     404 {object} errors.APIError "User not found"
// @Failure     500 {object} errors.APIError "Internal server error"
// @Router      /api/v1/users/{id} [get]
func (u *UserApi) GetById(c *gin.Context) {
	ctx := c.Request.Context()

	var params dto.GetUserByIdParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(
			errors.NewAPIError(
				http.StatusBadRequest,
				err.Error(),
			),
		)
		return
	}

	user, err := u.userService.GetById(ctx, params.ID)

	if err != nil {
		c.Error(err)
		return
	}

	response.WriteSuccess(c, user)
}
