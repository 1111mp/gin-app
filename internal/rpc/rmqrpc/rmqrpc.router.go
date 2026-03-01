package rmqrpc

import (
	"context"
	"fmt"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/internal/dto"
	post "github.com/1111mp/gin-app/internal/modules/post"
	"github.com/1111mp/gin-app/pkg/logger"
)

type RMQRPCRouter interface {
	GetPostById(ctx context.Context, req dto.GetPostByIdRequest) (*ent.PostEntity, error)
}

type rmqRPCRouterImpl struct {
	logger      logger.Logger
	postService post.PostService
}

func NewRMQRPCRouter(logger logger.Logger, postService post.PostService) RMQRPCRouter {
	return &rmqRPCRouterImpl{
		logger,
		postService,
	}
}

// GetPostById -.
func (c *rmqRPCRouterImpl) GetPostById(
	ctx context.Context,
	req dto.GetPostByIdRequest,
) (*ent.PostEntity, error) {
	post, err := c.postService.GetById(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("rmq_rpc - controller - v1 - GetPostById: %w", err)
	}

	return post, nil
}
