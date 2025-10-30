package api_v1

import (
	"github.com/1111mp/gin-app/config"
	api_service "github.com/1111mp/gin-app/internal/service/api"
)

// ApiGroup -.
type ApiGroup struct {
	UserApi        UserApiInter
	PostApi        PostApiInter
	AccessTokenApi AccessTokenApiInter
}

// NewApiGroup -.
func NewApiGroup(s *api_service.ServiceGroup, cfg config.ConfigInterface) *ApiGroup {
	return &ApiGroup{
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
