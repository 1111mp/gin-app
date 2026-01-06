package dto

// AuthLoginDto -.
type AuthLoginDto struct {
	Lang        string `form:"lang"`
	RedirectURL string `form:"redirect" binding:"required"`
}

// AuthCallbackDto -.
type AuthCallbackDto struct {
	State string `form:"state" binding:"required"`
	Code  string `form:"code" binding:"required"`
}
