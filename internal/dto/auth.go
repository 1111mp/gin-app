package dto

// AuthLoginDto -.
type AuthLoginDto struct {
	Client      string `form:"client" binding:"required" oneof:"google github"`
	Lang        string `form:"lang"`
	RedirectURL string `form:"redirect" binding:"required"`
}

// AuthLoginWithAccountDto -.
type AuthLoginWithAccountDto struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=12,max=32"`
	RedirectURL string `json:"redirect_url" binding:"omitempty,url"`
}

// AuthCallbackDto -.
type AuthCallbackDto struct {
	State string `form:"state" binding:"required"`
	Code  string `form:"code" binding:"required"`
}
