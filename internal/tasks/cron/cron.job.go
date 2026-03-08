package cron

import (
	"context"
	"time"

	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/1111mp/gin-app/pkg/rediskey"
	"github.com/go-redsync/redsync/v4"
)

type exampleJobImpl struct {
	logger logger.Logger
}

func NewExampleCronJob(
	logger logger.Logger,
	rs *redsync.Redsync,
	rdk rediskey.RedisKey,
) CronJob {
	return NewBaseCronJob(
		logger,
		rs,
		rdk,
		&exampleJobImpl{
			logger: logger,
		},
	)
}

func (j *exampleJobImpl) CronConfig() CronConfig {
	return CronConfig{
		Name:   "example_job",
		Spec:   "0 * * * * *", // schedule a cron job to run every minute
		Expiry: 60 * time.Second,
	}
}

func (j *exampleJobImpl) Do(ctx context.Context) {
	j.logger.Infof("cron job running")
	time.Sleep(20 * time.Second)
}
