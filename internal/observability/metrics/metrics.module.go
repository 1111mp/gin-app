package metrics

import (
	"context"

	"github.com/1111mp/gin-app/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"metrics",

	fx.Invoke(
		func(
			lc fx.Lifecycle,
			logger logger.Logger,
		) {
			lc.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						logger.Infof("app - Run - metrics - initializing")
						metricExporter, err := prometheus.New()
						if err != nil {
							logger.Errorf("app - Run - metrics - prometheus.New: %v", err)
							return err
						}
						mp := sdkmetric.NewMeterProvider(
							sdkmetric.WithReader(metricExporter),
						)
						otel.SetMeterProvider(mp)
						logger.Infof("app - Run - metrics initialized")
						return nil
					},
				},
			)
		},
	),
)
