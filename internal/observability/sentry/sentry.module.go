package sentry

import (
	"context"

	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/getsentry/sentry-go"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"sentry",

	fx.Invoke(
		func(
			lc fx.Lifecycle,
			cfg config.Config,
			logger logger.Logger,
		) {
			lc.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						err := sentry.Init(sentry.ClientOptions{
							Dsn:              cfg.Sentry().DSN,
							Environment:      cfg.App().Env,
							Release:          cfg.App().Version,
							AttachStacktrace: true,
							SampleRate:       1.0,
							SendDefaultPII:   false,
							BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
								if event.Request != nil {
									// remove sensitive headers
									delete(event.Request.Headers, "Cookie")
									delete(event.Request.Headers, "Authorization")
									delete(event.Request.Headers, "PRIVATE-TOKEN")
								}
								return event
							},
						})
						if err != nil {
							logger.Error("app - Run - sentry.Init failed", "error", err)
							return err
						}
						logger.Infof("app - Run - sentry initialized")
						return nil
					},
					OnStop: func(ctx context.Context) error {
						logger.Infof("app - Run - sentry shutting down")
						// use context-aware flush
						ok := sentry.FlushWithContext(ctx)
						if !ok {
							logger.Warnf("app - Run - sentry flush timeout")
						}
						logger.Infof("app - Run - sentry shutdown complete")
						return nil
					},
				},
			)
		},
	),
)
