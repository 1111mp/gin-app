package v1

import (
	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/internal/service"
)

// ApiGroup -.
type ApiGroup struct {
	UserApi UserApiInter
	PostApi PostApiInter
}

// NewApiGroup -.
func NewApiGroup(s *service.ServiceGroup, cfg config.ConfigInterface) *ApiGroup {
	return &ApiGroup{
		UserApi: &UserApi{
			cfg:         cfg,
			userService: s.UserService,
		},
		PostApi: &PostApi{
			postService: s.PostService,
		},
	}
}
