package v1

import (
	"net/http"

	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/internal/service"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/response"
	"github.com/gin-gonic/gin"
)

// UserApi -.
type UserApi struct {
	userService service.UserServiceInter
}

// GetById godoc
// @Summary     Get user by ID
// @Description Retrieve user information by given user ID
// @ID          getUserById
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Success     200 {object} response.UserAPIResponse "Successfully retrieved user"
// @Failure     400 {object} errors.APIError "Bad request (invalid ID)"
// @Failure     404 {object} errors.APIError "User not found"
// @Failure     500 {object} errors.APIError "Internal server error"
// @Router      /api/v1/user/{id} [get]
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
