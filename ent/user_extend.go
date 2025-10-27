package ent

// UserEntity -.
type UserEntity struct {
	ID         int           `json:"id,omitempty"`
	Name       string        `json:"name,omitempty"`
	Email      string        `json:"email,omitempty"`
	CreateTime string        `json:"createTime,omitempty"`
	UpdateTime string        `json:"updateTime,omitempty"`
	Posts      []*PostEntity `json:"posts"`
}

// IntoEntity converts ent User to UserEntity.
func (u *User) IntoEntity() *UserEntity {
	userEntity := &UserEntity{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		CreateTime: u.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: u.UpdateTime.Format("2006-01-02 15:04:05"),
	}

	if posts, err := u.Edges.PostsOrErr(); err == nil {
		userEntity.Posts = make([]*PostEntity, 0, len(posts))
		for _, post := range posts {
			userEntity.Posts = append(userEntity.Posts, post.IntoEntity())
		}
	}

	return userEntity
}
