package user

import (
	api_router "github.com/1111mp/gin-app/internal/router/api"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",

	fx.Provide(
		// user repository
		fx.Annotate(
			NewUserRepository,
			fx.As(new(UserRepository)),
		),
		// user service
		fx.Annotate(
			NewUserService,
			fx.As(new(UserService)),
		),
		// user controller
		NewUserController,
	),

	// register router
	fx.Invoke(func(
		router api_router.APIRouterParams,
		userController *UserController,
	) {
		// public routers
		{
			userGroup := router.Public.Group("/users")

			userGroup.POST("", userController.CreateOne)
		}

		// private routers
		{
			userGroup := router.Private.Group("/users")

			userGroup.GET("/:id", userController.GetById)
		}
	}),
)
