package post

import (
	"net/http"

	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/errors"
	"github.com/1111mp/gin-app/pkg/response"
	"github.com/gin-gonic/gin"
)

type PostController struct {
	postService PostService
}

func NewPostController(
	postService PostService,
) *PostController {
	return &PostController{postService}
}

// CreateOne godoc
// @Security 		APIAuth
// @Summary 		Create a new post
// @Description Creates a new post resource
// @ID          PostCreateOne
// @Tags 				Posts
// @Accept 			json
// @Produce 		json
// @Param 			data body dto.PostCreateOneDto true "Post data"
// @Success 		200 {object} response.PostAPIResponse "Post created successfully"
// @Failure     400 {object} errors.APIError "Bad request (invalid params)"
// @Failure     500 {object} errors.APIError "Internal server error"
// @Router 			/api/v1/posts [post]
func (p *PostController) CreateOne(c *gin.Context, userId int) {
	ctx := c.Request.Context()

	var dto dto.PostCreateOneDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(
			errors.NewAPIError(
				http.StatusBadRequest,
				err.Error(),
			),
		)
		return
	}

	post, err := p.postService.CreateOne(ctx, userId, dto)
	if err != nil {
		c.Error(err)
		return
	}

	response.WriteSuccess(c, post)
}

// GetById godoc
// @Security 		APIAuth
// @Summary			Get post by ID
// @Description	Retrieve post information by given post ID
// @ID					GetPostById
// @Tags				Posts
// @Accept      json
// @Produce     json
// @Param       id path int true "Post ID"
// @Success     200 {object} response.PostAPIResponse "Successfully retrieved post"
// @Failure     400 {object} errors.APIError "Bad request (invalid ID)"
// @Failure     404 {object} errors.APIError "Post not found"
// @Failure     500 {object} errors.APIError "Internal server error"
// @Router      /api/v1/posts/{id} [get]
func (p *PostController) GetById(c *gin.Context) {
	ctx := c.Request.Context()

	var params dto.PostGetByIdDto
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(
			errors.NewAPIError(
				http.StatusBadRequest,
				err.Error(),
			),
		)
		return
	}

	post, err := p.postService.GetById(ctx, params.ID)
	if err != nil {
		c.Error(err)
		return
	}

	response.WriteSuccess(c, post)
}
