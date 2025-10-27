package response

import "github.com/1111mp/gin-app/ent"

// UserAPIResponse -.
type UserAPIResponse struct {
	APIResponse[ent.UserEntity]
}
