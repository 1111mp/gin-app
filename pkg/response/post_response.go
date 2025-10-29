package response

import "github.com/1111mp/gin-app/ent"

// UserAPIResponse -.
type PostAPIResponse struct {
	APIResponse[ent.PostEntity]
}
