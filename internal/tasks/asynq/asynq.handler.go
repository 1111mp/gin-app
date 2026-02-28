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

type AsynqMuxHandler interface {
	HandleEmailDeliveryTask(ctx context.Context, t *asynqlib.Task) error
}

type asynqMuxHandlerImpl struct {
	logger logger.Logger
}

func NewAsynqMuxHandler(logger logger.Logger) AsynqMuxHandler {
	return &asynqMuxHandlerImpl{
		logger: logger,
	}
}

//---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input task.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
//---------------------------------------------------------------

// EmailDeliveryTask is an example task handler for processing email delivery tasks.
func (h *asynqMuxHandlerImpl) HandleEmailDeliveryTask(
	ctx context.Context,
	t *asynqlib.Task,
) error {
	var p dto.AsynqEmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		h.logger.Error(
			"[HandleEmailDeliveryTask] failed to unmarshal email delivery payload",
			"task_type", t.Type(),
			"payload", string(t.Payload()),
			"error", err,
		)
		return fmt.Errorf("[HandleEmailDeliveryTask] json.Unmarshal (task_type=%s): %v: %w", t.Type(), err, asynq.SkipRetry)
	}

	h.logger.Infof("Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	// TODO send email to user with given template.

	return nil
}
