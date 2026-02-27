package post

import (
	"context"
	"fmt"
	"net/http"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/errors"
)

type PostService interface {
	CreateOne(ctx context.Context, userId int, dto dto.PostCreateOneDto) (*ent.PostEntity, error)
	GetById(ctx context.Context, id int) (*ent.PostEntity, error)
}

type postServiceImpl struct {
	postRepository PostRepository
}

func NewPostService(
	postRepository PostRepository,
) PostService {
	return &postServiceImpl{postRepository}
}

// CreateOne -.
func (ps *postServiceImpl) CreateOne(ctx context.Context, userId int, dto dto.PostCreateOneDto) (*ent.PostEntity, error) {
	post, err := ps.postRepository.CreateOne(ctx, userId, dto)
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
func (ps *postServiceImpl) GetById(ctx context.Context, id int) (*ent.PostEntity, error) {
	post, err := ps.postRepository.GetById(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NewAPIError(
				http.StatusNotFound,
				fmt.Sprintf("post %d not found", id),
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
