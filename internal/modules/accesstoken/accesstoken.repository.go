package accesstoken

import (
	"context"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/ent/accesstoken"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/postgres"
)

type accessTokenRepositoryimpl struct {
	pg *postgres.Postgres
}

type AccessTokenRepository interface {
	CreateOne(ctx context.Context, userId int, dto dto.AccessTokenCreateOneDto) (*ent.AccessToken, error)
	GetByOwner(ctx context.Context, owner int) ([]*ent.AccessToken, error)
}

func NewAccessTokenRepository(pg *postgres.Postgres) AccessTokenRepository {
	return &accessTokenRepositoryimpl{pg}
}

// CreateOne -.
func (a *accessTokenRepositoryimpl) CreateOne(
	ctx context.Context,
	userId int,
	dto dto.AccessTokenCreateOneDto,
) (*ent.AccessToken, error) {
	return a.pg.Client.AccessToken.
		Create().
		SetName(dto.Name).
		SetOwner(dto.Owner).
		SetExpireTime(dto.ExpireTime).
		SetCreator(userId).
		Save(ctx)
}

// GetByOwner -.
func (a *accessTokenRepositoryimpl) GetByOwner(
	ctx context.Context,
	owner int,
) ([]*ent.AccessToken, error) {
	return a.pg.Client.AccessToken.
		Query().
		Where(
			accesstoken.OwnerEQ(owner),
		).
		All(ctx)
}
