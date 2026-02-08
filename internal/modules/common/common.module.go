package common

import (
	api_router "github.com/1111mp/gin-app/internal/router/api"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"common",

	fx.Provide(
		// common controller
		NewCommonController,
	),

	fx.Invoke(func(
		router api_router.APIRouterParams,
		commonController *CommonController,
	) {
		// public routers
		{
			commonGroup := router.Public.Group("/")

			commonGroup.GET("/healthz", commonController.Healthz)
		}

		// private routers
		{
			// commonGroup := router.Private.Group("/")
		}
	}),
)
