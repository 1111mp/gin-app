package api_v1

import (
	"net/http"

	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/1111mp/gin-app/pkg/oauth2"
	"github.com/1111mp/gin-app/pkg/oauth2/github"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthApiInter -.
type AuthApiInter interface {
	LoginHandler(c *gin.Context)
	CallbackHandler(c *gin.Context)
}

// AuthApi -.
type AuthApi struct {
	logger       logger.Interface
	githubOAuth2 github.OAuth2Inter
}

// LoginHandler godoc
// @Summary OAuth2 Login
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
	redirectURL := a.githubOAuth2.GetAuthURL(state)

	c.Redirect(http.StatusFound, redirectURL)
}

// CallbackHandler godoc
func (a *AuthApi) CallbackHandler(c *gin.Context) {
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

	token, err := a.githubOAuth2.GetToken(ctx, dto.Code)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	user, err := a.githubOAuth2.GetUser(ctx, token)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	// TODO store token and user info into DB
	a.logger.Info("user info", zap.Any("user", user))

	// response.WriteSuccess(c, user)
	if state.Redirect != "" {
		c.Redirect(http.StatusSeeOther, state.Redirect)
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}
