package tasks

import (
	"github.com/1111mp/gin-app/internal/tasks/asynq"
	"github.com/1111mp/gin-app/internal/tasks/cron"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"tasks",

	// asynq
	asynq.Module,
	// cron
	cron.Module,
)
