package asynq

import (
	"context"
	"fmt"

	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/pkg/logger"
	asynqlib "github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"asynq",

	fx.Provide(
		// asynq mux handlers for processing tasks
		NewAsynqMuxHandler,
		// asynq new task creator
		fx.Annotate(
			NewAsynqNewTask,
			// First parameter is logger.Logger, second is *asynqlib.Client
			// leave first param tag empty and tag the client parameter
			fx.ParamTags("", `name:"asynq_client"`),
			fx.ResultTags(`name:"asynq_new_task"`),
		),
		// asynq client
		fx.Annotate(
			func(lc fx.Lifecycle, cfg config.Config, logger logger.Logger) (*asynqlib.Client, error) {
				opt, err := asynqlib.ParseRedisURI(cfg.Redis().URL)
				if err != nil {
					logger.Fatal(fmt.Errorf("app - Run - asynq - client.ParseRedisURI: %w", err))
					return nil, err
				}

				client := asynqlib.NewClient(opt)
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						logger.Infof("app - Run - asynq - client shutting down")
						return client.Close()
					},
				})
				logger.Infof("app - Run - asynq - client initialized")
				return client, nil
			},
			fx.ResultTags(`name:"asynq_client"`),
		),
		// asynq server
		fx.Annotate(
			func(cfg config.Config, logger logger.Logger) (*asynqlib.Server, error) {
				opt, err := asynqlib.ParseRedisURI(cfg.Redis().URL)
				if err != nil {
					logger.Fatal(fmt.Errorf("app - Run - asynq - server.ParseRedisURI: %w", err))
					return nil, err
				}

				server := asynqlib.NewServer(opt, asynqlib.Config{
					// Specify how many concurrent workers to use
					Concurrency: 10,
					// Optionally specify multiple queues with different priority.
					Queues: map[string]int{
						"critical": 6,
						"default":  3,
						"low":      1,
					},
				})
				logger.Infof("app - Run - asynq - server initialized")
				return server, nil
			},
			fx.ResultTags(`name:"asynq_server"`),
		),
		// asynq scheduler
		// ! Scheduler is not used in this app, but you can easily add it back if you need to run periodic tasks.
		// ! If the task only needs to run periodically, robfig/cron is a better option.
		// ! It operates in-memory and avoids external dependencies such as Redis.
		// fx.Annotate(
		// 	func(lc fx.Lifecycle, cfg config.Config, logger logger.Logger) (*asynqlib.Scheduler, error) {
		// 		opt, err := asynqlib.ParseRedisURI(cfg.Redis().URL)
		// 		if err != nil {
		// 			logger.Fatal(fmt.Errorf("app - Run - asynq - server.ParseRedisURI: %w", err))
		// 			return nil, err
		// 		}

		// 		scheduler := asynqlib.NewScheduler(opt, nil)
		// 		logger.Infof("app - Run - asynq - scheduler initialized")
		// 		// Register periodic tasks
		// 		{
		// 			// if _, err := scheduler.Register("* * * * *", asynq.NewTask("task1", nil)); err != nil {
		// 			// 	log.Fatal(err)
		// 			// }
		// 		}

		// 		lc.Append(fx.Hook{
		// 			OnStart: func(ctx context.Context) error {
		// 				err := scheduler.Start()
		// 				if err != nil {
		// 					logger.Fatal(fmt.Errorf("app - Run - asynq - scheduler.Start: %w", err))
		// 				}
		// 				return nil
		// 			},
		// 			OnStop: func(ctx context.Context) error {
		// 				scheduler.Shutdown()
		// 				logger.Infof("app - Run - asynq - scheduler stopped")
		// 				return nil
		// 			},
		// 		})

		// 		return scheduler, nil
		// 	},
		// ),
		// asynq mux
		fx.Annotate(
			func(logger logger.Logger, h AsynqMuxHandler) *asynqlib.ServeMux {
				mux := asynqlib.NewServeMux()
				// Register task handlers
				{
					mux.HandleFunc(TypeEmailDelivery, h.HandleEmailDeliveryTask)
					mux.Handle(TypeImageResize, NewImageProcessor(logger))
				}

				logger.Infof("app - Run - asynq - mux initialized")
				return mux
			},
			fx.ResultTags(`name:"asynq_server_mux"`),
		),
	),

	fx.Invoke(
		fx.Annotate(
			func(
				lc fx.Lifecycle,
				sd fx.Shutdowner,
				logger logger.Logger,
				server *asynqlib.Server,
				mux *asynqlib.ServeMux,
			) {
				lc.Append(
					fx.Hook{
						OnStart: func(ctx context.Context) error {
							go func() {
								if err := server.Run(mux); err != nil {
									logger.Fatal(fmt.Errorf("app - Run - asynq - server.Run: %w", err))
									// ! It's better to just shut down the app and let Kubernetes restart it.
									_ = sd.Shutdown()
								}
							}()
							logger.Info("app - Run - asynq - server started")
							return nil
						},
						OnStop: func(ctx context.Context) error {
							done := make(chan struct{})
							go func() {
								server.Shutdown()
								close(done)
							}()

							select {
							case <-done:
								logger.Info("app - Run - asynq - server stopped")
							case <-ctx.Done():
								logger.Warn("app - Run - asynq - server shutdown timeout")
							}

							return nil
						},
					},
				)
			},
			// Use fx.ParamTags to specify the name of the server dependency
			fx.ParamTags("", "", "", `name:"asynq_server"`, `name:"asynq_server_mux"`),
		),
	),
)
