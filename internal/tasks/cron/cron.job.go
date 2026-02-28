package cron

import "github.com/1111mp/gin-app/pkg/logger"

type CronsJob interface {
	Run()
}

type cronJobsImpl struct {
	logger logger.Logger
}

func NewCronJobs(logger logger.Logger) CronsJob {
	return &cronJobsImpl{
		logger: logger,
	}
}

func (c *cronJobsImpl) Run() {
	c.logger.Infof("cron job running")
}
