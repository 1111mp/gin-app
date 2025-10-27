package ent

import "github.com/1111mp/gin-app/ent/post"

// PostEntity -.
type PostEntity struct {
	ID         int           `json:"id,omitempty"`
	Title      string        `json:"title,omitempty"`
	Content    string        `json:"content,omitempty"`
	Category   post.Category `json:"category,omitempty"`
	CreateTime string        `json:"createTime,omitempty"`
	UpdateTime string        `json:"updateTime,omitempty"`
}

// IntoEntity converts ent Post to PostEntity.
func (p *Post) IntoEntity() *PostEntity {
	return &PostEntity{
		ID:         p.ID,
		Title:      p.Title,
		Content:    p.Content,
		Category:   p.Category,
		CreateTime: p.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: p.UpdateTime.Format("2006-01-02 15:04:05"),
	}
}
