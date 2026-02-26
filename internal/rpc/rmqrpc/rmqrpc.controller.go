package rmqrpc

import (
	"context"
	"fmt"

	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/internal/dto"
	post "github.com/1111mp/gin-app/internal/modules/post"
)

type RMQRPCControllter interface {
	GetPostById(ctx context.Context, req dto.GetPostByIdRequest) (*ent.PostEntity, error)
}

type rmqRPCControllterImpl struct {
	postService post.PostService
}

func NewRMQRPCControllter(postService post.PostService) RMQRPCControllter {
	return &rmqRPCControllterImpl{postService}
}

// GetPostById -.
func (c *rmqRPCControllterImpl) GetPostById(
	ctx context.Context,
	req dto.GetPostByIdRequest,
) (*ent.PostEntity, error) {
	post, err := c.postService.GetById(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("rmq_rpc - controller - v1 - GetPostById: %w", err)
	}

	return post, nil
}
