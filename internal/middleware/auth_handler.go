package middleware

import (
	"net/http"
	"time"

	"github.com/1111mp/gin-app/ent/accesstoken"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/jwt"
	"github.com/1111mp/gin-app/pkg/postgres"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// APIAuthHandler -.
func APIAuthHandler(j jwt.JWT, name string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := trace.SpanFromContext(ctx.Request.Context())
		span.SetAttributes(
			attribute.String("auth.middleware", "api_auth"),
		)

		token, err := ctx.Cookie(name)
		if err != nil {
			span.SetAttributes(
				attribute.String("auth.middleware.result", "failed"),
				attribute.String("auth.middleware.reason", "missing_cookie"),
			)
			ctx.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}

		claims, err := j.ParseToken(token)
		if err != nil {
			span.SetAttributes(
				attribute.String("auth.middleware.result", "failed"),
				attribute.String("auth.middleware.reason", "invalid_token"),
			)
			ctx.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}

		if claims == nil || claims.UserId == 0 {
			span.SetAttributes(
				attribute.String("auth.result", "failed"),
				attribute.String("auth.reason", "invalid_claims"),
			)
			ctx.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}

		span.SetAttributes(
			attribute.String("auth.result", "success"),
		)

		ctx.Set("userId", claims.UserId)
		ctx.Next()
	}
}

// OpenAPIAuthHandler -.
func OpenAPIAuthHandler(pg *postgres.Postgres) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := ctx.GetHeader("PRIVATE-TOKEN")
		if accessToken == "" {
			ctx.AbortWithError(
				http.StatusUnauthorized,
				errors.NewAPIError(
					http.StatusUnauthorized,
					"Authentication error: The request header did not include a 'PRIVATE-TOKEN'.",
				),
			)
			return
		}

		at, err := pg.Client.AccessToken.
			Query().
			Where(accesstoken.ValueEQ(accessToken)).
			Only(ctx.Request.Context())
		if err != nil {
			ctx.AbortWithError(
				http.StatusUnauthorized,
				errors.NewAPIError(http.StatusUnauthorized, "Invalid token."),
			)
			return
		}

		if at.ExpireTime > 0 && at.ExpireTime < time.Now().Unix() {
			ctx.AbortWithError(
				http.StatusUnauthorized,
				errors.NewAPIError(http.StatusUnauthorized, "The token has expired."),
			)
			return
		}

		ctx.Set("userId", at.Owner)
		ctx.Next()
	}
}
