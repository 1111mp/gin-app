package state

import (
	"github.com/1111mp/gin-app/pkg/postgres"
	"github.com/1111mp/gin-app/pkg/redis"
)

// AppState -.
type AppState struct {
	PG    *postgres.Postgres
	Redis *redis.Redis
}
