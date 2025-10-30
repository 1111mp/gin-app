package openapi_v1

import (
	openapi_service "github.com/1111mp/gin-app/internal/service/open-api"
)

// ApiGroup -.
type ApiGroup struct {
}

// NewApiGroup -.
func NewApiGroup(s *openapi_service.ServiceGroup) *ApiGroup {
	return &ApiGroup{}
}
