package middleware

import (
	"errors"
	"net/http"

	appErrors "github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/gin-gonic/gin"
)

// ErrorHandler captures errors and returns a consistent JSON error response.
func ErrorHandler(l logger.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			var repErr *appErrors.RepositoryError
			if errors.As(err, &repErr) {
				l.Errorf("[repository error]: %v", repErr)

				ctx.JSON(
					http.StatusInternalServerError,
					gin.H{
						"code":    http.StatusInternalServerError,
						"message": repErr.Message,
					},
				)
				return
			}

			var apiErr *appErrors.APIError
			if errors.As(err, &apiErr) {
				ctx.JSON(
					apiErr.Code,
					gin.H{
						"code":    apiErr.Code,
						"message": apiErr.Message,
					},
				)
				return
			}

			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"message": err.Error(),
				},
			)
		}
	}
}
