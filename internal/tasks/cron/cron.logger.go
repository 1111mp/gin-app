package cron

import (
	"github.com/1111mp/gin-app/pkg/logger"
	cronlib "github.com/robfig/cron/v3"
)

type CronLogger struct {
	log logger.Logger
}

var _ cronlib.Logger = (*CronLogger)(nil)

func NewCronLogger(log logger.Logger) cronlib.Logger {
	return &CronLogger{log: log}
}

func (c *CronLogger) Info(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) > 0 {
		c.log.Infow(msg, keysAndValues...)
		return
	}
	c.log.Info(msg)
}

func (c *CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) > 0 {
		kv := append([]interface{}{"error", err}, keysAndValues...)
		c.log.Errorw(msg, kv...)
		return
	}
	c.log.Errorf("%s: %v", msg, err)
}
