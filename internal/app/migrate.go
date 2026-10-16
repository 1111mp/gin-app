//go:build migrate

package app

import (
	"errors"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 10
	_defaultTimeout  = time.Second
)

func init() {
	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		log.Fatalf("migrate: environment variable not declared: PG_URL")
	}

	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		log.Fatalf("migrate: invalid PG_URL: %v", err)
	}
	query := parsedURL.Query()
	if query.Get("sslmode") == "" {
		query.Set("sslmode", "disable")
	}
	parsedURL.RawQuery = query.Encode()
	databaseURL = parsedURL.String()

	var (
		attempts = _defaultAttempts
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://ent/migrate/migrations", databaseURL)
		if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}
