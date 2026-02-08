package accesstoken

import (
	"context"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/errors"
)

type accessTokenServiceImpl struct {
	accessTokenRepository AccessTokenRepository
}

type AccessTokenService interface {
	CreateOne(ctx context.Context, userId int, dto dto.AccessTokenCreateOneDto) (*ent.AccessTokenEntity, error)
	GetByOwner(ctx context.Context, owner int) ([]*ent.AccessTokenEntity, error)
}

func NewAccessTokenService(accessTokenRepository AccessTokenRepository) AccessTokenService {
	return &accessTokenServiceImpl{accessTokenRepository}
}

// CreateOne -.
func (a *accessTokenServiceImpl) CreateOne(ctx context.Context, userId int, dto dto.AccessTokenCreateOneDto) (*ent.AccessTokenEntity, error) {
	accessToken, err := a.accessTokenRepository.CreateOne(ctx, userId, dto)
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

// GetByOwner -.
func (a *accessTokenServiceImpl) GetByOwner(ctx context.Context, owner int) ([]*ent.AccessTokenEntity, error) {
	accessTokens, err := a.accessTokenRepository.GetByOwner(ctx, owner)
	if err != nil {
		return nil, errors.WrapAPIError(
			errors.ErrInternalServerError,
			errors.NewRepositoryError(
				err.Error(),
				err,
			),
		)
	}

	var entities []*ent.AccessTokenEntity
	for _, accessToken := range accessTokens {
		entities = append(entities, accessToken.IntoEntity())
	}

	return entities, nil
}
