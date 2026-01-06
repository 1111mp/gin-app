package repository

import (
	"context"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/ent/accesstoken"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/state"
)

// AccessTokenRepositoryInter -.
type AccessTokenRepositoryInter interface {
	CreateOne(ctx context.Context, userId int, dto dto.AccessTokenCreateOneDto) (*ent.AccessToken, error)
	GetByOwner(ctx context.Context, owner int) ([]*ent.AccessToken, error)
}

// AccessTokenRepository -.
type AccessTokenRepository struct {
	appState *state.AppState
}

// CreateOne -.
func (a *AccessTokenRepository) CreateOne(
	ctx context.Context,
	userId int,
	dto dto.AccessTokenCreateOneDto,
) (*ent.AccessToken, error) {
	return a.appState.PG.Client.AccessToken.
		Create().
		SetName(dto.Name).
		SetOwner(dto.Owner).
		SetExpireTime(dto.ExpireTime).
		SetCreator(userId).
		Save(ctx)
}

// GetByOwner -.
func (a *AccessTokenRepository) GetByOwner(
	ctx context.Context,
	owner int,
) ([]*ent.AccessToken, error) {
	return a.appState.PG.Client.AccessToken.
		Query().
		Where(
			accesstoken.OwnerEQ(owner),
		).
		All(ctx)
}
