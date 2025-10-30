package dto

// AccessTokenCreateOneDto -.
type AccessTokenCreateOneDto struct {
	Name       string `json:"name" binding:"required"`
	Owner      int    `json:"owner" binding:"required"`
	ExpireTime int64  `json:"expireTime" binding:"required"`
}
