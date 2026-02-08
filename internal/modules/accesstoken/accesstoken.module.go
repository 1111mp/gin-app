package accesstoken

import (
	api_router "github.com/1111mp/gin-app/internal/router/api"
	"github.com/1111mp/gin-app/pkg/utils"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"accessToken",

	fx.Provide(
		// accessToken repository
		fx.Annotate(
			NewAccessTokenRepository,
			fx.As(new(AccessTokenRepository)),
		),
		// accessToken service
		fx.Annotate(
			NewAccessTokenService,
			fx.As(new(AccessTokenService)),
		),
		// accessToken controller
		NewAccessTokenController,
	),

	// register router
	fx.Invoke(func(
		router api_router.APIRouterParams,
		accessTokenController *AccessTokenController,
	) {
		// public routers
		{
			// accessTokenGroup := router.Public.Group("/access-tokens")
		}

		// private routers
		{
			accessTokenGroup := router.Private.Group("/access-tokens")

			accessTokenGroup.POST("", utils.HandlerWithUser(accessTokenController.CreateOne))

			accessTokenGroup.GET("", utils.HandlerWithUser(accessTokenController.GetSelfTokens))
			accessTokenGroup.GET("/:owner", accessTokenController.GetTokensByOwner)
		}
	}),
)
