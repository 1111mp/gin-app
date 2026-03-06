package middleware

import (
	"errors"
	"net/http"

	appErrors "github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ErrorHandler captures errors and returns a consistent JSON error response.
func ErrorHandler(logger logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors.Last().Err
		span := trace.SpanFromContext(ctx.Request.Context())

		var repErr *appErrors.RepositoryError
		if errors.As(err, &repErr) {
			span.RecordError(err)
			span.SetStatus(codes.Error, repErr.Message)

			logger.Errorw(
				"log from middleware error handler",
				"error", repErr.Message,
				"request_id", requestid.Get(ctx),
				"method", ctx.Request.Method,
				"path", ctx.Request.URL.Path,
				"route", ctx.FullPath(),
				"handler", ctx.HandlerName(),
			)

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
			if apiErr.Code >= 500 {
				span.RecordError(err)
				span.SetStatus(codes.Error, apiErr.Message)
			}

			ctx.JSON(
				apiErr.Code,
				gin.H{
					"code":    apiErr.Code,
					"message": apiErr.Message,
				},
			)
			return
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			},
		)

	}
}
