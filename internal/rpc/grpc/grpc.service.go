package grpc

import (
	"context"
	"fmt"

	v1 "github.com/1111mp/gin-app/docs/proto/v1"
	"github.com/1111mp/gin-app/internal/modules/post"
	"github.com/1111mp/gin-app/pkg/logger"
)

type GRPCSvervice interface {
	v1.PostServer
}

type grpcServiceImpl struct {
	v1.PostServer

	logger      logger.Logger
	postService post.PostService
}

func NewGRPCSvervice(logger logger.Logger, postService post.PostService) GRPCSvervice {
	return &grpcServiceImpl{
		logger:      logger,
		postService: postService,
	}
}

func (c *grpcServiceImpl) GetPostById(ctx context.Context, req *v1.GetPostByIdRequest) (*v1.GetPostByIdResponse, error) {
	post, err := c.postService.GetById(ctx, int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("grpc - controller - v1 - GetPostById: %w", err)
	}

	return &v1.GetPostByIdResponse{
		Id:         int64(post.ID),
		Title:      post.Title,
		Content:    post.Content,
		Category:   v1.Category(v1.Category_value[string(post.Category)]),
		CreateTime: post.CreateTime,
		UpdateTime: post.UpdateTime,
	}, nil
}
