package tasks

import (
	"github.com/1111mp/gin-app/internal/tasks/asynq"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"tasks",

	// asynq
	asynq.Module,
)
