package cron

import (
	"context"

	"github.com/1111mp/gin-app/pkg/logger"
	cronlib "github.com/robfig/cron/v3"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"schedule",

	fx.Provide(
		// cron jobs
		fx.Annotate(
			NewCronJobs,
			fx.ResultTags(`name:"cron_jobs"`),
		),
		// corn
		fx.Annotate(
			func(
				logger logger.Logger,
			) *cronlib.Cron {
				cronLogger := NewCronLogger(logger)
				cron := cronlib.New(
					cronlib.WithSeconds(),
					cronlib.WithChain(
						cronlib.Recover(cronLogger),
						cronlib.SkipIfStillRunning(cronLogger),
					),
				)
				logger.Infof("app - Run - cron - initialized")
				return cron
			},
			fx.ResultTags(`name:"cron"`),
		),
	),
	// cron jobs registration
	fx.Invoke(
		fx.Annotate(
			func(
				logger logger.Logger,
				cron *cronlib.Cron,
				cronJobs CronsJob,
			) {
				// register cron jobs here
				{
					// schedule a cron job to run every minute
					cron.AddFunc("0 * * * * *", cronJobs.Run)
				}

				logger.Infof("app - Run - cron - jobs registered")
			},
			fx.ParamTags("", `name:"cron"`, `name:"cron_jobs"`),
		),
	),

	fx.Invoke(
		fx.Annotate(
			func(
				lc fx.Lifecycle,
				sd fx.Shutdowner,
				logger logger.Logger,
				cron *cronlib.Cron,
			) {
				lc.Append(
					fx.Hook{
						OnStart: func(ctx context.Context) error {
							cron.Start()
							logger.Infof("app - Run - cron - started")
							return nil
						},
						OnStop: func(ctx context.Context) error {
							ctx = cron.Stop()
							<-ctx.Done()
							logger.Infof("app - Run - cron - stopped")
							return nil
						},
					},
				)
			},
			fx.ParamTags("", "", "", `name:"cron"`),
		),
	),
)
