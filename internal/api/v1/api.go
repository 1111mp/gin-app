package api_v1

import (
	"github.com/1111mp/gin-app/config"
	api_service "github.com/1111mp/gin-app/internal/service/api"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/1111mp/gin-app/pkg/oauth2/github"
	"github.com/1111mp/gin-app/pkg/oauth2/google"
)

// ApiGroup -.
type ApiGroup struct {
	AuthApi        AuthApiInter
	UserApi        UserApiInter
	PostApi        PostApiInter
	AccessTokenApi AccessTokenApiInter
}

// NewApiGroup -.
func NewApiGroup(
	s *api_service.ServiceGroup,
	cfg config.ConfigInterface,
	logger logger.Interface,
) *ApiGroup {
	// setup github oauth2 client
	githubClient := github.Setup(
		github.ClientID(cfg.Github().ClientID),
		github.ClientSecret(cfg.Github().ClientSecret),
		github.RedirectURL(cfg.Github().RedirectURL),
	)

	// setup google oauth2 client
	googleClient := google.Setup(
		google.ClientID(cfg.Google().ClientID),
		google.ClientSecret(cfg.Google().ClientSecret),
		google.RedirectURL(cfg.Google().RedirectURL),
	)

	return &ApiGroup{
		&AuthApi{
			logger,
			googleClient,
			githubClient,
		},
		&UserApi{
			cfg:         cfg,
			userService: s.UserService,
		},
		&PostApi{
			postService: s.PostService,
		},
		&AccessTokenApi{
			accessTokenService: s.AccessTokenService,
		},
	}
}
