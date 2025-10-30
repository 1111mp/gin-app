package response

import "github.com/1111mp/gin-app/ent"

// AccessTokenAPIResponse -.
type AccessTokenAPIResponse struct {
	APIResponse[ent.AccessTokenEntity]
}
