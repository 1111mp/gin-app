package modules

import (
	"github.com/1111mp/gin-app/internal/modules/accesstoken"
	"github.com/1111mp/gin-app/internal/modules/auth"
	"github.com/1111mp/gin-app/internal/modules/common"
	"github.com/1111mp/gin-app/internal/modules/post"
	"github.com/1111mp/gin-app/internal/modules/user"
	"go.uber.org/fx"
)

var APIModule = fx.Module(
	"api",

	// auth module
	auth.Module,
	// user module
	user.Module,
	// post module
	post.Module,
	// access-token module
	accesstoken.Module,
	// common module
	common.Module,
)
