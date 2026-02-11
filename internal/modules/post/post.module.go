package post

import (
	api_router "github.com/1111mp/gin-app/internal/router/api"
	"github.com/1111mp/gin-app/pkg/utils"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"post",

	fx.Provide(
		// post repository
		NewPostRepository,
		// post service
		NewPostService,
		// post controller
		NewPostController,
	),

	// register router
	fx.Invoke(func(
		router api_router.APIRouterParams,
		postController *PostController,
	) {
		// public routers
		{
			// postGroup := router.Public.Group("/posts")
		}

		// private routers
		{
			postGroup := router.Private.Group("/posts")

			postGroup.POST("", utils.HandlerWithUser(postController.CreateOne))
			postGroup.GET("/:id", postController.GetById)
		}
	}),
)
