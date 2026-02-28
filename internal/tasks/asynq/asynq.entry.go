package asynq

import (
	"time"

	"github.com/1111mp/gin-app/internal/dto"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/goccy/go-json"
	asynqlib "github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeEmailDelivery = "email:deliver"
)

type AsynqNewTask interface {
	NewEmailDeliveryTask(emailDeliveryPayload dto.AsynqEmailDeliveryPayload) (*asynqlib.TaskInfo, error)
}

type asynqNewTaskImpl struct {
	logger logger.Logger
	client *asynqlib.Client
}

func NewAsynqNewTask(logger logger.Logger, client *asynqlib.Client) AsynqNewTask {
	return &asynqNewTaskImpl{
		logger: logger,
		client: client,
	}
}

// NewEmailDeliveryTask creates a new email delivery task with the given payload.
func (t *asynqNewTaskImpl) NewEmailDeliveryTask(
	emailDeliveryPayload dto.AsynqEmailDeliveryPayload,
) (*asynqlib.TaskInfo, error) {
	payload, err := json.Marshal(emailDeliveryPayload)
	if err != nil {
		return nil, err
	}
	//task options can be passed to NewTask, which can be overridden at enqueue time.
	task := asynqlib.NewTask(TypeEmailDelivery, payload, asynqlib.MaxRetry(5), asynqlib.Timeout(10*time.Minute))
	// Enqueue the task to be processed.
	return t.client.Enqueue(task)
}

// NewImageResizeTask creates a new image resize task with the given payload.
func (t *asynqNewTaskImpl) NewImageResizeTask(
	imageResizePayload dto.AsynqImageResizePayload,
) (*asynqlib.TaskInfo, error) {
	payload, err := json.Marshal(imageResizePayload)
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	task := asynqlib.NewTask(TypeImageResize, payload, asynqlib.MaxRetry(5), asynqlib.Timeout(10*time.Minute))
	// Enqueue the task to be processed.
	return t.client.Enqueue(task, asynqlib.MaxRetry(5), asynqlib.Timeout(5*time.Minute))
}
