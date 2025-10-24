package repository

import "github.com/1111mp/gin-app/pkg/postgres"

//go:generate mockgen -source=interfaces.go -destination=../service/mocks_user_test.go -package=service_test

type UserRepository interface {
	CreateOne()
}

// UserRepository -.
type userRepository struct {
	pg *postgres.Postgres
}

// CreateOne -.
func (u *userRepository) CreateOne() {}
