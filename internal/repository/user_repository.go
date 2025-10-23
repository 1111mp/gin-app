package repository

import "github.com/1111mp/gin-app/pkg/postgres"

// UserRepository -.
type UserRepository struct {
	pg *postgres.Postgres
}

// CreateOne -.
func (u *UserRepository) CreateOne() {}
