package repository

import (
	"context"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/ent/user"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/postgres"
)

//go:generate mockgen -source=interfaces.go -destination=../service/mocks_user_test.go -package=service_test

// UserRepository -.
type UserRepositoryInter interface {
	CreateOne(ctx context.Context, dto dto.CreateOneUserDto) (*ent.User, error)
	GetById(ctx context.Context, id int) (*ent.User, error)
}

// UserRepository -.
type UserRepository struct {
	pg *postgres.Postgres
}

// CreateOne -.
func (u *UserRepository) CreateOne(
	ctx context.Context,
	dto dto.CreateOneUserDto,
) (*ent.User, error) {
	return u.pg.Client.User.
		Create().
		SetName(dto.Name).
		SetEmail(dto.Email).
		SetPassword(dto.Password).
		Save(ctx)
}

// GetById -.
func (u *UserRepository) GetById(
	ctx context.Context,
	id int,
) (*ent.User, error) {
	return u.pg.Client.User.
		Query().
		WithPosts().
		Where(
			user.IDEQ(id),
		).
		Only(ctx)
}
