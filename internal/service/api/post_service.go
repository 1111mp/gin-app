package api_service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/internal/repository"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/logger"
)

// PostServiceInter-.
type PostServiceInter interface {
	CreateOne(ctx context.Context, userId int, dto dto.PostCreateOneDto) (*ent.PostEntity, error)
	GetById(ctx context.Context, id int) (*ent.PostEntity, error)
}

// PostService -.
type PostService struct {
	l   logger.Interface
	rep repository.PostRepositoryInter
}

// CreateOne -.
func (p *PostService) CreateOne(ctx context.Context, userId int, dto dto.PostCreateOneDto) (*ent.PostEntity, error) {
	post, err := p.rep.CreateOne(ctx, userId, dto)
	if err != nil {
		return nil, errors.WrapAPIError(
			errors.ErrInternalServerError,
			errors.NewRepositoryError(
				err.Error(),
				err,
			),
		)
	}

	return post.IntoEntity(), nil
}

// GetById -.
func (p *PostService) GetById(ctx context.Context, id int) (*ent.PostEntity, error) {
	post, err := p.rep.GetById(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.WrapAPIError(
				errors.NewAPIError(
					http.StatusNotFound,
					fmt.Sprintf("post %d not found", id),
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

	return post.IntoEntity(), nil
}
