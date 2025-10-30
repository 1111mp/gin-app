package response

import "github.com/1111mp/gin-app/ent"

// PostAPIResponse -.
type PostAPIResponse struct {
	APIResponse[ent.PostEntity]
}
