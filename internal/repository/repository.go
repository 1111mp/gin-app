package repository

import "github.com/1111mp/gin-app/pkg/postgres"

// RepositoryGroup -.
type RepositoryGroup struct {
	UserRepository UserRepository
}

// NewRepositoryGroup -.
func NewRepositoryGroup(pg *postgres.Postgres) *RepositoryGroup {
	return &RepositoryGroup{
		UserRepository: &userRepository{
			pg,
		},
	}
}
