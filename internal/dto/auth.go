package dto

// AuthLoginDto -.
type AuthLoginDto struct {
	Client      string `form:"client" binding:"required" oneof:"google github"`
	Lang        string `form:"lang"`
	RedirectURL string `form:"redirect" binding:"required"`
}

// AuthCallbackDto -.
type AuthCallbackDto struct {
	State string `form:"state" binding:"required"`
	Code  string `form:"code" binding:"required"`
}
