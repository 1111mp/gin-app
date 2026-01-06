package repository

import (
	"context"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/ent/user"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/state"
)

//go:generate mockgen -source=user_repository.go -destination=../service/api/mocks_user_repo_test.go -package=api_service_test

// UserRepositoryInter -.
type UserRepositoryInter interface {
	CreateOne(ctx context.Context, dto dto.UserCreateOneDto) (*ent.User, error)
	GetById(ctx context.Context, id int) (*ent.User, error)
}

// UserRepository -.
type UserRepository struct {
	appState *state.AppState
}

// CreateOne -.
func (u *UserRepository) CreateOne(
	ctx context.Context,
	dto dto.UserCreateOneDto,
) (*ent.User, error) {
	return u.appState.PG.Client.User.
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
	return u.appState.PG.Client.User.
		Query().
		WithPosts().
		Where(
			user.IDEQ(id),
		).
		Only(ctx)
}
