package api_service

import (
	"context"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/internal/repository"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/logger"
)

// AccessTokenServiceInter -.
type AccessTokenServiceInter interface {
	CreateOne(ctx context.Context, userId int, dto dto.AccessTokenCreateOneDto) (*ent.AccessTokenEntity, error)
}

// AccessTokenService -.
type AccessTokenService struct {
	l   logger.Interface
	rep repository.AccessTokenRepositoryInter
}

// CreateOne -.
func (a *AccessTokenService) CreateOne(ctx context.Context, userId int, dto dto.AccessTokenCreateOneDto) (*ent.AccessTokenEntity, error) {
	accessToken, err := a.rep.CreateOne(ctx, userId, dto)
	if err != nil {
		return nil, errors.WrapAPIError(
			errors.ErrInternalServerError,
			errors.NewRepositoryError(
				err.Error(),
				err,
			),
		)
	}

	return accessToken.IntoEntity(), nil
}
