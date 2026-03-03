package observability

import (
	"github.com/1111mp/gin-app/internal/observability/metrics"
	"github.com/1111mp/gin-app/internal/observability/sentry"
	"github.com/1111mp/gin-app/internal/observability/tracing"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"observability",

	// tracing
	tracing.Module,
	// metrics
	metrics.Module,
	// sentry
	sentry.Module,
)
