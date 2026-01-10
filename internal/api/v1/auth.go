package api_v1

import (
	"net/http"

	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/1111mp/gin-app/pkg/oauth2"
	"github.com/1111mp/gin-app/pkg/oauth2/github"
	"github.com/1111mp/gin-app/pkg/oauth2/google"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthApiInter -.
type AuthApiInter interface {
	LoginHandler(c *gin.Context)
	GoogleCallbackHandler(c *gin.Context)
	GithubCallbackHandler(c *gin.Context)
}

// AuthApi -.
type AuthApi struct {
	logger       logger.Interface
	googleClient google.ClientInter
	githubClient github.ClientInter
}

// LoginHandler godoc
// @Summary     Redirect to GitHub OAuth2 login page
// @Description Redirects the user to GitHub's OAuth2 login page for authentication
// @ID          AuthLogin
// @Tags        Auth
// @Param       redirect query string true  "URL to redirect after successful login"
// @Param       lang     query string false "Preferred language"
// @Success     302 "Redirects to GitHub OAuth2 login page"
// @Failure     400 {object} errors.APIError "Bad request (invalid params)"
// @Router      /api/v1/auth/login [get]
func (a *AuthApi) LoginHandler(c *gin.Context) {
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

	state, err := oauth2.MakeState(dto.RedirectURL, dto.Lang)
	if err != nil {
		c.Error(err)
		return
	}

	var redirectURL string
	switch dto.Client {
	case "google":
		redirectURL = a.googleClient.GetAuthURL(state)
	case "github":
		redirectURL = a.githubClient.GetAuthURL(state)
	default:
		// unsupported client
		c.Error(errors.NewAPIError(http.StatusBadRequest, "Unsupported oauth2 client"))
		return
	}

	c.Redirect(http.StatusFound, redirectURL)
}

// GoogleCallbackHandler -.
func (a *AuthApi) GoogleCallbackHandler(c *gin.Context) {
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

	// TODO check state to avoid CSRF attacks
	state, err := oauth2.ParseState(dto.State)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	token, err := a.googleClient.GetToken(ctx, dto.Code)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	user, err := a.googleClient.GetUser(ctx, token)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	// TODO store token and user info into DB
	a.logger.Info("user info", zap.Any("user", user))

	if state.Redirect != "" {
		c.Redirect(http.StatusSeeOther, state.Redirect)
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}

// GithubCallbackHandler -.
func (a *AuthApi) GithubCallbackHandler(c *gin.Context) {
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

	// TODO check state to avoid CSRF attacks
	state, err := oauth2.ParseState(dto.State)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	token, err := a.githubClient.GetToken(ctx, dto.Code)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	user, err := a.githubClient.GetUser(ctx, token)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	// TODO store token and user info into DB
	a.logger.Info("user info", zap.Any("user", user))

	if state.Redirect != "" {
		c.Redirect(http.StatusSeeOther, state.Redirect)
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}
