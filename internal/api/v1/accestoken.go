package api_v1

import (
	"net/http"

	"github.com/1111mp/gin-app/internal/dto"
	api_service "github.com/1111mp/gin-app/internal/service/api"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/response"
	"github.com/gin-gonic/gin"
)

// AccessTokenApiInter -.
type AccessTokenApiInter interface {
	CreateOne(c *gin.Context, userId int)
}

// AccessTokenApi -.
type AccessTokenApi struct {
	accessTokenService api_service.AccessTokenServiceInter
}

// CreateOne godoc
// @Summary 		Create a new access token
// @Description Creates a new access token resource
// @ID          AccessTokenCreateOne
// @Tags 				AccessTokens
// @Accept 			json
// @Produce 		json
// @Param 			data body dto.AccessTokenCreateOneDto true "AccessToken data"
// @Success 		200 {object} response.AccessTokenAPIResponse "Post created successfully"
// @Failure     400 {object} errors.APIError "Bad request (invalid params)"
// @Failure     500 {object} errors.APIError "Internal server error"
// @Router 			/api/v1/access-tokens [post]
func (a *AccessTokenApi) CreateOne(c *gin.Context, userId int) {
	ctx := c.Request.Context()

	var dto dto.AccessTokenCreateOneDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(
			errors.NewAPIError(
				http.StatusBadRequest,
				err.Error(),
			),
		)
		return
	}

	accessToken, err := a.accessTokenService.CreateOne(ctx, userId, dto)
	if err != nil {
		c.Error(err)
		return
	}

	response.WriteSuccess(c, accessToken)
}
