package common

import (
	api_router "github.com/1111mp/gin-app/internal/router/api"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
			commonGroup.GET("/metrics", gin.WrapH(promhttp.Handler()))
		}

		// private routers
		{
			// commonGroup := router.Private.Group("/")
		}
	}),
)
