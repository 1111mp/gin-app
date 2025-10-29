package repository

import (
	"context"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/ent/post"
	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/postgres"
)

// PostRepositoryInter -.
type PostRepositoryInter interface {
	CreateOne(ctx context.Context, userId int, dto dto.PostCreateOneDto) (*ent.Post, error)
	GetById(ctx context.Context, id int) (*ent.Post, error)
}

// PostRepository -.
type PostRepository struct {
	pg *postgres.Postgres
}

// CreateOne -.
func (p *PostRepository) CreateOne(
	ctx context.Context,
	userId int,
	dto dto.PostCreateOneDto,
) (*ent.Post, error) {
	return p.pg.Client.Post.
		Create().
		SetOwnerID(userId).
		SetTitle(dto.Title).
		SetContent(dto.Content).
		SetCategory(dto.Category).
		Save(ctx)
}

// GetById -.
func (p *PostRepository) GetById(
	ctx context.Context,
	id int,
) (*ent.Post, error) {
	return p.pg.Client.Post.
		Query().
		Where(
			post.IDEQ(id),
		).
		Only(ctx)
}
