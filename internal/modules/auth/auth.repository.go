package auth

import (
	"context"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/ent/user"
	"github.com/1111mp/gin-app/pkg/postgres"
	"go.opentelemetry.io/otel"
)

var repTracer = otel.Tracer("AuthRepository")

type AuthRepository interface {
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
}

type authRepositoryImpl struct {
	pg *postgres.Postgres
}

func NewAuthRepository(pg *postgres.Postgres) AuthRepository {
	return &authRepositoryImpl{pg}
}

// GetByEmail -.
func (a *authRepositoryImpl) GetByEmail(
	ctx context.Context,
	email string,
) (*ent.User, error) {
	ctx, span := repTracer.Start(ctx, "AuthRepository.GetByEmail")
	defer span.End()

	user, err := a.pg.Client.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}
