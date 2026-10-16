package cron

import (
	"context"
	"time"

	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/1111mp/gin-app/pkg/rediskey"
	"github.com/go-redsync/redsync/v4"
	"go.opentelemetry.io/otel"
)

const cronUnlockTimeout = 500 * time.Millisecond

var tracer = otel.Tracer("cron")

type CronJob interface {
	Run()
	Spec() string
	Cancel()
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
	ctx    context.Context
	cancel context.CancelFunc
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
	ctx, cancel := context.WithCancel(context.Background())

	return &baseCronJobImpl{
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
		rs:     rs,
		key:    rdk.Key("cron", cfg.Name),
		spec:   cfg.Spec,
		expiry: cfg.Expiry,
		job:    job,
	}
}

func (b *baseCronJobImpl) Run() {
	ctx, span := tracer.Start(b.ctx, b.key)
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
		unlockCtx, cancel := context.WithTimeout(context.Background(), cronUnlockTimeout)
		defer cancel()

		_, err := mutex.UnlockContext(unlockCtx)
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

func (b *baseCronJobImpl) Cancel() { b.cancel() }
