package openapi_router

import (
	"github.com/1111mp/gin-app/internal/middleware"
	"github.com/1111mp/gin-app/pkg/postgres"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type OpenAPIRouter struct {
	fx.Out

	Private *gin.RouterGroup `name:"open-api:private"`
}

type OpenAPIRouterParams struct {
	fx.In

	Private *gin.RouterGroup `name:"open-api:private"`
}

func NewOpenAPIRouter(
	app *gin.Engine,
	pg *postgres.Postgres,
) OpenAPIRouter {
	private := app.Group("/open-api/v1")

	private.Use(middleware.OpenAPIAuthHandler(pg))

	return OpenAPIRouter{
		Private: private,
	}
}
