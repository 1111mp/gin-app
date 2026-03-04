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

	fx.Provide(
		fx.Annotate(
			func(logger logger.Logger) (*sdkmetric.MeterProvider, error) {
				metricExporter, err := prometheus.New()
				if err != nil {
					logger.Errorf("app - Run - metrics - prometheus.New: %v", err)
					return nil, err
				}
				mp := sdkmetric.NewMeterProvider(
					sdkmetric.WithReader(metricExporter),
				)
				logger.Infof("app - Run - metrics - Provider initialized")
				return mp, nil
			},
			fx.ResultTags(`name:"metrics_provider"`),
		),
	),

	fx.Invoke(
		fx.Annotate(
			func(
				lc fx.Lifecycle,
				logger logger.Logger,
				mp *sdkmetric.MeterProvider,
			) {
				lc.Append(
					fx.Hook{
						OnStart: func(ctx context.Context) error {
							otel.SetMeterProvider(mp)
							logger.Infof("app - Run - metrics - set global meter provider")
							return nil
						},
					},
				)
			},
			fx.ParamTags("", "", `name:"metrics_provider"`),
		),
	),
)
