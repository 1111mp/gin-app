package tracing

import (
	"context"
	"time"

	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"tracing",

	fx.Provide(
		// tracing filter
		fx.Annotate(
			NewTraceFilter,
			fx.ResultTags(`name:"tracing_filter"`),
		),
		// tracing provider
		fx.Annotate(
			func(cfg config.Config, logger logger.Logger) (*sdktrace.TracerProvider, error) {
				// use background context for initialization
				ctx := context.Background()
				// 1️⃣ OTLP gRPC client
				client := otlptracegrpc.NewClient(
					// Tempo OTLP gRPC
					otlptracegrpc.WithEndpoint("localhost:4319"),
					// Insecure connection (no TLS)
					otlptracegrpc.WithInsecure(),
					// Timeout for the connection
					otlptracegrpc.WithTimeout(3*time.Second),
				)
				// 2️⃣ exporter
				exporter, err := otlptrace.New(ctx, client)
				if err != nil {
					logger.Errorf("app - Run - tracing - otlptrace.New: %v", err)
					return nil, err
				}
				// 3️⃣ resource
				res, err := resource.New(
					ctx,
					resource.WithAttributes(
						semconv.ServiceName(cfg.App().Name),
						semconv.ServiceVersion(cfg.App().Version),
						semconv.DeploymentEnvironmentName(cfg.App().Env),
					),
				)
				if err != nil {
					logger.Errorf("app - Run - tracing - resource.New: %v", err)
					return nil, err
				}
				// 4️⃣ tracer provider
				tp := sdktrace.NewTracerProvider(
					sdktrace.WithBatcher(exporter),
					sdktrace.WithResource(res),
					sdktrace.WithSampler(
						sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0)),
					),
				)
				logger.Infof("app - Run - tracing provider initialized")
				return tp, nil
			},
			fx.ResultTags(`name:"tracing_provider"`),
		),
	),

	fx.Invoke(
		fx.Annotate(
			func(
				lc fx.Lifecycle,
				logger logger.Logger,
				tp *sdktrace.TracerProvider,
			) {
				lc.Append(
					fx.Hook{
						OnStart: func(ctx context.Context) error {
							// set global tracer provider
							otel.SetTracerProvider(tp)
							// set global propagator to tracecontext (the default is no-op).
							otel.SetTextMapPropagator(
								propagation.NewCompositeTextMapPropagator(
									propagation.TraceContext{},
									propagation.Baggage{},
								),
							)
							logger.Infof("app - Run - tracing - set global tracer provider and propagator")
							return nil
						},
						OnStop: func(ctx context.Context) error {
							logger.Infof("app - Run - tracing - shutting down")
							err := tp.Shutdown(ctx)
							if err != nil {
								logger.Errorf("app - Run - tracing - tp.Shutdown: %v", err)
								return err
							}
							logger.Infof("app - Run - tracing - shutdown complete")
							return nil
						},
					},
				)
			},
			fx.ParamTags("", "", `name:"tracing_provider"`),
		),
	),
)
