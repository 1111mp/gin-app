package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/1111mp/gin-app/ent"
	_ "github.com/1111mp/gin-app/ent/runtime"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool   *pgxpool.Pool
	Client *ent.Client
}

// New -.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)
	// ✅ Set up OpenTelemetry tracing for pgx
	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	for pg.connAttempts > 0 {
		pool, poolErr := pgxpool.NewWithConfig(context.Background(), poolConfig)
		if poolErr != nil {
			err = poolErr
			log.Printf("Postgres is trying to connect, attempts left: %d, err: %v", pg.connAttempts, err)

			time.Sleep(pg.connTimeout)
			pg.connAttempts--

			continue
		}

		// ✅ Ping to check connection
		ctx, cancel := context.WithTimeout(context.Background(), pg.connTimeout)
		err = pool.Ping(ctx)
		cancel()

		if err == nil {
			pg.Pool = pool
			break
		}

		pool.Close()

		log.Printf("Postgres is trying to connect, attempts left: %d, err: %v", pg.connAttempts, err)

		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	// ✅ Record database stats with OpenTelemetry
	if err = otelpgx.RecordStats(pg.Pool); err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - unable to record database stats: %w", err)
	}

	// ✅ Initialize Ent client
	pg.Client = NewEntClient(pg.Pool)

	return pg, nil
}

// NewEntClient -.
func NewEntClient(pool *pgxpool.Pool) *ent.Client {
	db := stdlib.OpenDBFromPool(pool)
	drv := entsql.OpenDB(dialect.Postgres, db)

	return ent.NewClient(ent.Driver(drv))
}

// Close -.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
