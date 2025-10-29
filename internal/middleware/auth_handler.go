package middleware

import (
	"net/http"

	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// AuthHandler -.
func AuthHandler(j jwt.JWTManagerInterface, name string) gin.HandlerFunc {
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
