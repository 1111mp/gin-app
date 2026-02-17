package auth

import (
	"net/http"

	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	cfg         config.Config
	authService AuthService
}

func NewAuthController(cfg config.Config, authService AuthService) *AuthController {
	return &AuthController{cfg, authService}
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

// LoginWithAccount godoc
// @Summary     Login with account
// @Description Login with account
// @ID          AuthLoginWithAccount
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       data body dto.AuthLoginWithAccountDto true "Login data"
// @Success     302 {string} string "Redirect to login page"
// @Failure     400 {object} errors.APIError "Bad request (invalid params)"
// @Failure     500 {object} errors.APIError "Internal server error"
// @Router      /api/v1/auth/login-with-account [post]
func (a *AuthController) LoginWithAccount(c *gin.Context) {
	ctx := c.Request.Context()

	var dto dto.AuthLoginWithAccountDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(
			errors.NewAPIError(
				http.StatusBadRequest,
				err.Error(),
			),
		)
		return
	}

	token, err := a.authService.LoginWithAccount(ctx, dto)
	if err != nil {
		c.Error(err)
		return
	}

	// set cookie
	c.SetCookie(a.cfg.HTTP().CookieName, token, 3600, "/", "", true, true)

	if dto.RedirectURL != "" {
		c.Redirect(http.StatusFound, dto.RedirectURL)
	} else {
		c.Redirect(http.StatusFound, "/")
	}
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
