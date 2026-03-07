package cron

import (
	"context"
	"time"

	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/1111mp/gin-app/pkg/rediskey"
	"github.com/go-redsync/redsync/v4"
)

type cronJobImpl struct {
	*baseCronJobImpl

	spec string
}

func NewCronJobs(logger logger.Logger, rs *redsync.Redsync, rdk rediskey.RedisKey) CronJob {
	job := &cronJobImpl{
		// schedule a cron job to run every minute
		spec: "0 * * * * *",
	}
	base := NewBaseCronJob(logger, rs, rdk.Key("cron", "example_job"), 30*time.Second, job)
	job.baseCronJobImpl = base
	return job
}

func (j *cronJobImpl) Do(ctx context.Context) {
	j.logger.Infof("cron job running")
	time.Sleep(20 * time.Second)
}

func (j *cronJobImpl) Spec() string {
	return j.spec
}
