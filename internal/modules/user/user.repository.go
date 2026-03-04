package user

import (
	"context"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/ent/user"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/postgres"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

//go:generate mockgen -source=user_repository.go -destination=./mocks_user_repo_test.go -package=user_test

var repTracer = otel.Tracer("UserRepository")

// UserRepository -.
type UserRepository interface {
	CreateOne(ctx context.Context, dto dto.UserCreateOneDto) (*ent.User, error)
	GetById(ctx context.Context, id int) (*ent.User, error)
}

type userRepositoryImpl struct {
	pg *postgres.Postgres
}

func NewUserRepository(pg *postgres.Postgres) UserRepository {
	return &userRepositoryImpl{pg}
}

// CreateOne -.
func (u *userRepositoryImpl) CreateOne(
	ctx context.Context,
	dto dto.UserCreateOneDto,
) (*ent.User, error) {
	return u.pg.Client.User.
		Create().
		SetName(dto.Name).
		SetEmail(dto.Email).
		SetPassword(dto.Password).
		Save(ctx)
}

// GetById -.
func (u *userRepositoryImpl) GetById(
	ctx context.Context,
	id int,
) (*ent.User, error) {
	ctx, span := repTracer.Start(ctx, "UserRepository.GetById")
	defer span.End()

	span.SetAttributes(
		attribute.Bool("has_user_id", true),
	)

	user, err := u.pg.Client.User.
		Query().
		WithPosts().
		Where(
			user.IDEQ(id),
		).
		Only(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}
