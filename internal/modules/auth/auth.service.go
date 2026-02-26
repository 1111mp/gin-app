package auth

import (
	"context"
	"net/http"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/jwt"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/1111mp/gin-app/pkg/oauth2"
	"github.com/1111mp/gin-app/pkg/oauth2/github"
	"github.com/1111mp/gin-app/pkg/oauth2/google"
	"go.uber.org/zap"
)

type AuthService interface {
	Login(dto dto.AuthLoginDto) (string, error)
	LoginWithAccount(ctx context.Context, dto dto.AuthLoginWithAccountDto) (string, error)
	GithubCallback(ctx context.Context, dto dto.AuthCallbackDto) (*oauth2.State, error)
	GoogleCallback(ctx context.Context, dto dto.AuthCallbackDto) (*oauth2.State, error)
}

type authServiceImpl struct {
	logger         logger.Logger
	jwt            jwt.JWT
	githubClient   github.Client
	googleClient   google.Client
	authRepository AuthRepository
}

func NewAuthService(
	logger logger.Logger,
	jwt jwt.JWT,
	githubClient github.Client,
	googleClient google.Client,
	authRepository AuthRepository,
) AuthService {
	return &authServiceImpl{logger, jwt, githubClient, googleClient, authRepository}
}

// Login -.
func (a *authServiceImpl) Login(dto dto.AuthLoginDto) (string, error) {
	state, err := oauth2.MakeState(dto.RedirectURL, dto.Lang)
	if err != nil {
		return "", errors.WrapAPIError(
			errors.ErrInternalServerError,
			err,
		)
	}

	var redirectURL string
	switch dto.Client {
	case "google":
		redirectURL = a.googleClient.GetAuthURL(state)
	case "github":
		redirectURL = a.githubClient.GetAuthURL(state)
	default:
		// unsupported client
		return "", errors.NewAPIError(http.StatusBadRequest, "Unsupported oauth2 client")
	}

	return redirectURL, nil
}

// LoginWithAccount -.
func (a *authServiceImpl) LoginWithAccount(ctx context.Context, dto dto.AuthLoginWithAccountDto) (string, error) {
	user, err := a.authRepository.GetByEmail(ctx, dto.Email)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", errors.ErrUnauthorized
		}

		return "", errors.WrapAPIError(
			errors.ErrInternalServerError,
			errors.NewRepositoryError(
				err.Error(),
				err,
			),
		)
	}

	// compare password
	if err := user.ComparePassword(dto.Password); err != nil {
		return "", errors.ErrUnauthorized
	}

	token, err := a.jwt.GenerateToken(user.ID)
	if err != nil {
		return "", errors.WrapAPIError(
			errors.ErrInternalServerError,
			errors.NewRepositoryError(
				err.Error(),
				err,
			),
		)
	}

	return token, nil
}

// GithubCallback -.
func (a *authServiceImpl) GithubCallback(ctx context.Context, dto dto.AuthCallbackDto) (*oauth2.State, error) {
	// TODO check state to avoid CSRF attacks
	state, err := oauth2.ParseState(dto.State)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	token, err := a.githubClient.GetToken(ctx, dto.Code)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	user, err := a.githubClient.GetUser(ctx, token)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	// TODO store token and user info into DB
	a.logger.Info("user info", zap.Any("user", user))

	return state, nil
}

// GoogleCallback -.
func (a *authServiceImpl) GoogleCallback(ctx context.Context, dto dto.AuthCallbackDto) (*oauth2.State, error) {
	// TODO check state to avoid CSRF attacks
	state, err := oauth2.ParseState(dto.State)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	token, err := a.googleClient.GetToken(ctx, dto.Code)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	user, err := a.googleClient.GetUser(ctx, token)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	// TODO store token and user info into DB
	a.logger.Info("user info", zap.Any("user", user))

	return state, nil
}
