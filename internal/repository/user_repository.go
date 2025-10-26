package repository

import (
	"context"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/ent/user"
	"github.com/1111mp/gin-app/pkg/postgres"
)

//go:generate mockgen -source=interfaces.go -destination=../service/mocks_user_test.go -package=service_test

// UserRepository -.
type UserRepository interface {
	CreateOne()
	GetById(ctx context.Context, id int) (*ent.User, error)
}

// UserRepository -.
type userRepository struct {
	pg *postgres.Postgres
}

// CreateOne -.
func (u *userRepository) CreateOne() {}

// GetById -.
func (u *userRepository) GetById(ctx context.Context, id int) (*ent.User, error) {
	user, err := u.pg.Client.User.
		Query().
		WithPosts().
		Where(
			user.IDEQ(id),
		).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}
