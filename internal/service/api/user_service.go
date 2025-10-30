package api_service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/internal/repository"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/jwt"
	"github.com/1111mp/gin-app/pkg/logger"
)

// UserServiceInter -.
type UserServiceInter interface {
	CreateOne(ctx context.Context, dto dto.UserCreateOneDto) (*ent.UserEntity, string, error)
	GetById(ctx context.Context, id int) (*ent.UserEntity, error)
}

// UserService -.
type UserService struct {
	l   logger.Interface
	rep repository.UserRepositoryInter
	jwt jwt.JWTManagerInterface
}

// NewUserService -.
func NewUserService(l logger.Interface,
	rep repository.UserRepositoryInter,
	jwt jwt.JWTManagerInterface,
) *UserService {
	return &UserService{
		l:   l,
		rep: rep,
		jwt: jwt,
	}
}

// CreateUser -.
func (u *UserService) CreateOne(ctx context.Context, dto dto.UserCreateOneDto) (*ent.UserEntity, string, error) {
	user, err := u.rep.CreateOne(ctx, dto)
	if err != nil {
		return nil, "", errors.WrapAPIError(
			errors.ErrInternalServerError,
			errors.NewRepositoryError(
				err.Error(),
				err,
			),
		)
	}

	token, err := u.jwt.GenerateToken(user.ID)
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
func (u *UserService) GetById(ctx context.Context, id int) (*ent.UserEntity, error) {
	user, err := u.rep.GetById(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.WrapAPIError(
				errors.NewAPIError(
					http.StatusNotFound,
					fmt.Sprintf("user %d not found", id),
				),
				err,
			)
		}

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
