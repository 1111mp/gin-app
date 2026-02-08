package auth

import (
	"net/http"

	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService AuthService
}

func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{authService}
}

// Login godoc
// @Summary     Redirect to GitHub OAuth2 login page
// @Description Redirects the user to GitHub's OAuth2 login page for authentication
// @ID          AuthLogin
// @Tags        Auth
// @Param       redirect query string true  "URL to redirect after successful login"
// @Param       lang     query string false "Preferred language"
// @Success     302 "Redirects to GitHub OAuth2 login page"
// @Failure     400 {object} errors.APIError "Bad request (invalid params)"
// @Router      /api/v1/auth/login [get]
func (a *AuthController) Login(c *gin.Context) {
	var dto dto.AuthLoginDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.Error(
			errors.NewAPIError(
				http.StatusBadRequest,
				err.Error(),
			),
		)
		return
	}

	redirectURL, err := a.authService.Login(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusFound, redirectURL)
}

// GithubCallback -.
func (a *AuthController) GithubCallback(c *gin.Context) {
	ctx := c.Request.Context()

	var dto dto.AuthCallbackDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.Error(
			errors.NewAPIError(
				http.StatusBadRequest,
				err.Error(),
			),
		)
		return
	}

	state, err := a.authService.GithubCallback(ctx, dto)
	if err != nil {
		c.Error(err)
		return
	}

	if state.Redirect != "" {
		c.Redirect(http.StatusSeeOther, state.Redirect)
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}

// GoogleCallback -.
func (a *AuthController) GoogleCallback(c *gin.Context) {
	ctx := c.Request.Context()

	var dto dto.AuthCallbackDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.Error(
			errors.NewAPIError(
				http.StatusBadRequest,
				err.Error(),
			),
		)
		return
	}

	state, err := a.authService.GoogleCallback(ctx, dto)
	if err != nil {
		c.Error(err)
		return
	}

	if state.Redirect != "" {
		c.Redirect(http.StatusSeeOther, state.Redirect)
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}
