package cron

import (
	"context"
	"time"

	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/go-redsync/redsync/v4"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("cron")

type CronJob interface {
	Run()
	Spec() string
}

type BaseCronJob interface {
	Do(ctx context.Context)
}

type baseCronJobImpl struct {
	logger logger.Logger
	rs     *redsync.Redsync

	key    string
	expiry time.Duration

	self BaseCronJob
}

func NewBaseCronJob(
	logger logger.Logger,
	rs *redsync.Redsync,
	key string,
	expiry time.Duration,
	job BaseCronJob,
) *baseCronJobImpl {
	return &baseCronJobImpl{
		logger: logger,
		rs:     rs,
		key:    key,
		expiry: expiry,
		self:   job,
	}
}

func (b *baseCronJobImpl) Run() {
	ctx, span := tracer.Start(context.Background(), b.key)
	defer span.End()

	mutex := b.rs.NewMutex(
		b.key,
		redsync.WithExpiry(b.expiry),
		redsync.WithTries(1),
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

	if b.self != nil {
		b.self.Do(ctx)
	}

	b.logger.Infof("cron job %s completed", b.key)
}
