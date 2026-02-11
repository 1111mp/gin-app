package auth

import (
	api_router "github.com/1111mp/gin-app/internal/router/api"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"auth",

	fx.Provide(
		// auth service
		NewAuthService,
		// auth controller
		NewAuthController,
	),

	// register router
	fx.Invoke(func(
		router api_router.APIRouterParams,
		authController *AuthController,
	) {
		// public routers
		{
			// authGroup := router.Public.Group("/auth")
		}

		// private routers
		{
			authGroup := router.Private.Group("/auth")

			authGroup.GET("/login", authController.Login)
			authGroup.GET("/callback/google", authController.GoogleCallback)
			authGroup.GET("/callback/github", authController.GithubCallback)
		}
	}),
)
