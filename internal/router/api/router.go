package api_router

import (
	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/internal/middleware"
	"github.com/1111mp/gin-app/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type APIRouter struct {
	fx.Out

	Public  *gin.RouterGroup `name:"api:public"`
	Private *gin.RouterGroup `name:"api:private"`
}

type APIRouterParams struct {
	fx.In

	Public  *gin.RouterGroup `name:"api:public"`
	Private *gin.RouterGroup `name:"api:private"`
}

func NewAPIRouter(
	cfg config.Config,
	jwt jwt.JWT,
	app *gin.Engine,
) APIRouter {
	public := app.Group("/api/v1")
	private := public.Group("/")

	private.Use(middleware.APIAuthHandler(jwt, cfg.HTTP().CookieName))

	return APIRouter{
		Public:  public,
		Private: private,
	}
}
