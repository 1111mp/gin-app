package asynq

import (
	"context"
	"fmt"

	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	asynqlib "github.com/hibiken/asynq"
)

const (
	TypeImageResize = "image:resize"
)

// ImageProcessor implements asynqlib.Handler interface.
type ImageProcessor interface {
	ProcessTask(context.Context, *asynqlib.Task) error
}

type imageProcessorImpl struct {
	logger logger.Logger
}

func NewImageProcessor(logger logger.Logger) ImageProcessor {
	return &imageProcessorImpl{
		logger: logger,
	}
}

// ImageProcessor implements asynqlib.Handler interface.
func (processor *imageProcessorImpl) ProcessTask(
	ctx context.Context,
	t *asynqlib.Task,
) error {
	var p dto.AsynqImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		processor.logger.Error(
			"[HandleImageResizeTask] failed to unmarshal image resize payload",
			"task_type", t.Type(),
			"payload", string(t.Payload()),
			"error", err,
		)
		return fmt.Errorf("[HandleImageResizeTask] json.Unmarshal (task_type=%s): %v: %w", t.Type(), err, asynq.SkipRetry)
	}

	processor.logger.Infof("Resizing image: src=%s", p.SourceURL)
	// TODO Image resizing code ...

	return nil
}
