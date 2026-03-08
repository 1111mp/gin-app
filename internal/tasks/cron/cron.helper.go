package cron

import (
	"context"
	"time"

	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/1111mp/gin-app/pkg/rediskey"
	"github.com/go-redsync/redsync/v4"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("cron")

type CronJob interface {
	Run()
	Spec() string
}

type CronConfig struct {
	Name   string
	Spec   string
	Expiry time.Duration
}

type baseCronJob interface {
	Do(ctx context.Context)
	CronConfig() CronConfig
}

type baseCronJobImpl struct {
	logger logger.Logger
	rs     *redsync.Redsync

	key    string
	spec   string
	expiry time.Duration
	job    baseCronJob
}

func NewBaseCronJob(
	logger logger.Logger,
	rs *redsync.Redsync,
	rdk rediskey.RedisKey,
	job baseCronJob,
) CronJob {
	cfg := job.CronConfig()

	return &baseCronJobImpl{
		logger: logger,
		rs:     rs,
		key:    rdk.Key("cron", cfg.Name),
		spec:   cfg.Spec,
		expiry: cfg.Expiry,
		job:    job,
	}
}

func (b *baseCronJobImpl) Run() {
	ctx, span := tracer.Start(context.Background(), b.key)
	defer span.End()

	mutex := b.rs.NewMutex(
		b.key,
		redsync.WithExpiry(b.expiry),
		redsync.WithTries(12),
	)

	if err := mutex.LockContext(ctx); err != nil {
		b.logger.Infof("cron job %s already running", b.key)
		return
	}

	defer func() {
		_, err := mutex.UnlockContext(ctx)
		if err != nil {
			b.logger.Errorf("cron job %s unlock failed: %v", b.key, err)
		}
	}()

	b.logger.Infof("cron job %s started", b.key)

	if b.job != nil {
		b.job.Do(ctx)
	}

	b.logger.Infof("cron job %s completed", b.key)
}

func (b *baseCronJobImpl) Spec() string {
	return b.spec
}
