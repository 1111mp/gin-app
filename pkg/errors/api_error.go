package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError - 接口层面返回的错误
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{Code: code, Message: message}
}

func WrapAPIError(base *APIError, err error) *APIError {
	if err == nil {
		return base
	}
	base.Err = err
	return base
}

func FromError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return &APIError{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
		Err:     err,
	}
}

var (
	ErrBadRequest          = NewAPIError(http.StatusBadRequest, "Bad Request")
	ErrUnauthorized        = NewAPIError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden           = NewAPIError(http.StatusForbidden, "Forbidden")
	ErrNotFound            = NewAPIError(http.StatusNotFound, "Resource Not Found")
	ErrConflict            = NewAPIError(http.StatusConflict, "Conflict")
	ErrInternalServerError = NewAPIError(http.StatusInternalServerError, "Internal Server Error")
)
