package cron

import (
	"context"
	"testing"
	"time"

	"github.com/1111mp/gin-app/pkg/logger"
)

func TestExampleJobHonorsCancellation(t *testing.T) {
	t.Parallel()

	job := &exampleJobImpl{logger: logger.New("", "error")}
	ctx, cancel := context.WithCancel(t.Context())
	finished := make(chan struct{})

	go func() { job.Do(ctx); close(finished) }()

	cancel()

	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("job ignored cancellation")
	}
}
