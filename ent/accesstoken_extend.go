package ent

// AccessTokenEntity -.
type AccessTokenEntity struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Owner      int    `json:"owner,omitempty"`
	Creator    int    `json:"creator,omitempty"`
	CreateTime string `json:"createTime,omitempty"`
	UpdateTime string `json:"updateTime,omitempty"`
}

// IntoEntity converts ent AccessToken to AccessTokenEntity.
func (a *AccessToken) IntoEntity() *AccessTokenEntity {
	return &AccessTokenEntity{
		ID:         a.ID,
		Name:       a.Name,
		Owner:      a.Owner,
		Creator:    a.Creator,
		CreateTime: a.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: a.UpdateTime.Format("2006-01-02 15:04:05"),
	}
}
