package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/jwt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var srvTracer = otel.Tracer("UserService")

type UserService interface {
	CreateOne(ctx context.Context, dto dto.UserCreateOneDto) (*ent.UserEntity, string, error)
	GetById(ctx context.Context, id int) (*ent.UserEntity, error)
}

type userServiceImpl struct {
	jwt            jwt.JWT
	userRepository UserRepository
}

func NewUserService(jwt jwt.JWT, userRepository UserRepository) UserService {
	return &userServiceImpl{jwt, userRepository}
}

// CreateUser -.
func (us *userServiceImpl) CreateOne(ctx context.Context, dto dto.UserCreateOneDto) (*ent.UserEntity, string, error) {
	user, err := us.userRepository.CreateOne(ctx, dto)
	if err != nil {
		return nil, "", errors.WrapAPIError(
			errors.ErrInternalServerError,
			errors.NewRepositoryError(
				err.Error(),
				err,
			),
		)
	}

	token, err := us.jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, "", errors.WrapAPIError(
			errors.ErrInternalServerError,
			errors.NewRepositoryError(
				err.Error(),
				err,
			),
		)
	}

	return user.IntoEntity(), token, nil
}

// GetById -.
func (us *userServiceImpl) GetById(ctx context.Context, id int) (*ent.UserEntity, error) {
	ctx, span := srvTracer.Start(ctx, "UserService.GetById")
	defer span.End()

	span.SetAttributes(
		attribute.Bool("has_user_id", true),
	)

	user, err := us.userRepository.GetById(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			span.SetAttributes(
				attribute.Bool("user.not_found", true),
			)
			return nil, errors.NewAPIError(
				http.StatusNotFound,
				fmt.Sprintf("user %d not found", id),
			)
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WrapAPIError(
			errors.ErrInternalServerError,
			errors.NewRepositoryError(
				err.Error(),
				err,
			),
		)
	}

	return user.IntoEntity(), nil
}
