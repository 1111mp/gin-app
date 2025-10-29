package utils

import (
	"net/http"

	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/gin-gonic/gin"
)

// MustGetUser retrieves the current request's userId from the Gin context.
// Returns ErrUnauthorized if the userId does not exist or is not an int.
func MustGetUser(c *gin.Context) (int, error) {
	val, exists := c.Get("userId")
	if !exists {
		return 0, errors.ErrUnauthorized
	}

	id, ok := val.(int)
	if !ok {
		return 0, errors.ErrUnauthorized
	}

	return id, nil
}

// HandlerWithUser is a handler wrapper that injects userId into the handler.
// The handler signature must be func(c *gin.Context, userId int).
// If userId is missing or invalid, the request is aborted with a 401 Unauthorized response.
func HandlerWithUser(handler func(c *gin.Context, userId int)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := MustGetUser(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		handler(ctx, userId)
	}
}
