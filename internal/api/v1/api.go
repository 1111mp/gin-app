package v1

import (
	"github.com/1111mp/gin-app/internal/service"
)

// ApiGroup -.
type ApiGroup struct {
	UserApi *UserApi
}

// NewApiGroup -.
func NewApiGroup(s *service.ServiceGroup) *ApiGroup {
	return &ApiGroup{
		UserApi: &UserApi{
			userService: s.UserService,
		},
	}
}
