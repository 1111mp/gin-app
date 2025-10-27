package dto

type GetUserByIdParams struct {
	ID int `uri:"id" binding:"required,min=1"`
}
