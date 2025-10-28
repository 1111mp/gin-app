package dto

// CreateOneUserDto -.
type CreateOneUserDto struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=12,max=32"`
}

// GetUserByIdParams -.
type GetUserByIdParams struct {
	ID int `uri:"id" binding:"required,min=1"`
}
