package dto

import "github.com/1111mp/gin-app/ent/post"

// PostCreateOneDto -.
type PostCreateOneDto struct {
	Title    string        `json:"title" binding:"required"`
	Content  string        `json:"content" binding:"required"`
	Category post.Category `json:"category" binding:"required,oneof=Feed Story"`
}

// GetPostByIdDto -.
type PostGetByIdDto struct {
	ID int `uri:"id" binding:"required,min=1"`
}
