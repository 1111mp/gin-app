package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse -.
type APIResponse[T any] struct {
	Code    int    `json:"code" example:"200"`
	Payload *T     `json:"payload,omitempty"`
	Message string `json:"message" example:"success"`
}

// Success -.
func Success[T any](payload *T) *APIResponse[T] {
	return &APIResponse[T]{
		Code:    http.StatusOK,
		Payload: payload,
		Message: "success",
	}
}

// WriteSuccess -.
func WriteSuccess[T any](c *gin.Context, payload *T) {
	c.JSON(http.StatusOK, Success(payload))
}
