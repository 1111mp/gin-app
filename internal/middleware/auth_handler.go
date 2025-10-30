package middleware

import (
	"net/http"
	"time"

	"github.com/1111mp/gin-app/ent/accesstoken"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/jwt"
	"github.com/1111mp/gin-app/pkg/postgres"
	"github.com/gin-gonic/gin"
)

// APIAuthHandler -.
func APIAuthHandler(j jwt.JWTManagerInterface, name string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie(name)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}

		claims, err := j.ParseToken(token)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}

		if claims == nil || claims.UserId == 0 {
			ctx.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}

		ctx.Set("userId", claims.UserId)
		ctx.Next()
	}
}

// OpenAPIAuthHandler -.
func OpenAPIAuthHandler(pg *postgres.Postgres) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := ctx.GetHeader("PRIVATE-TOKEN")
		if accessToken == "" {
			ctx.AbortWithError(http.StatusUnauthorized, errors.NewAPIError(http.StatusUnauthorized, "Authentication error: The request header did not include a 'PRIVATE-TOKEN'."))
			return
		}

		at, err := pg.Client.AccessToken.
			Query().
			Where(accesstoken.ValueEQ(accessToken)).
			Only(ctx.Request.Context())
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.NewAPIError(http.StatusUnauthorized, "Invalid token."))
			return
		}

		if at.ExpireTime > 0 && at.ExpireTime < time.Now().Unix() {
			ctx.AbortWithError(http.StatusUnauthorized, errors.NewAPIError(http.StatusUnauthorized, "The token has expired."))
			return
		}

		ctx.Set("userId", at.Owner)
		ctx.Next()
	}
}
