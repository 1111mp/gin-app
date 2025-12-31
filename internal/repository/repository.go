package repository

import "github.com/1111mp/gin-app/pkg/state"

// RepositoryGroup -.
type RepositoryGroup struct {
	UserRepository        UserRepositoryInter
	PostRepository        PostRepositoryInter
	AccessTokenRepository AccessTokenRepositoryInter
}

// NewRepositoryGroup -.
func NewRepositoryGroup(appState *state.AppState) *RepositoryGroup {
	return &RepositoryGroup{
		&UserRepository{
			appState,
		},
		&PostRepository{
			appState,
		},
		&AccessTokenRepository{
			appState,
		},
	}
}
