package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/1111mp/gin-app/config"
	v1 "github.com/1111mp/gin-app/docs/proto/v1"
	grpcclient "github.com/1111mp/gin-app/pkg/grpc/client"
	grpcserver "github.com/1111mp/gin-app/pkg/grpc/server"
	"github.com/1111mp/gin-app/pkg/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc/reflection"
)

const grpcStopTimeout = 5 * time.Second

var Module = fx.Module(
	"grpc",

	fx.Provide(
		// grpc service
		NewGRPCSvervice,
		// grpc server
		fx.Annotate(
			func(cfg config.Config, logger logger.Logger) *grpcserver.Server {
				return grpcserver.New(logger, grpcserver.Port(cfg.GRPC().Port))
			},
		),
		// grpc client
		fx.Annotate(
			func(cfg config.Config, logger logger.Logger) (*grpcclient.Client, error) {
				return grpcclient.New(logger, grpcclient.Port(cfg.GRPC().Port))
			},
			fx.OnStop(func(client *grpcclient.Client) error {
				return client.Shutdown()
			}),
		),
	),
	// register router
	fx.Invoke(
		func(
			cfg config.Config,
			grpcServer *grpcserver.Server,
			grpcService GRPCSvervice,
		) {
			{
				v1.RegisterPostServer(grpcServer.App, grpcService)
			}

			if cfg.App().Env == "development" {
				reflection.Register(grpcServer.App)
			}
		},
	),

	fx.Invoke(
		func(
			lc fx.Lifecycle,
			sd fx.Shutdowner,
			logger logger.Logger,
			grpcServer *grpcserver.Server,
		) {
			lc.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						grpcServer.Start()

						go func() {
							if err, ok := <-grpcServer.Notify(); ok {
								logger.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
								if err := sd.Shutdown(fx.ExitCode(1)); err != nil {
									logger.Error(err)
								}
							}
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						ctx, cancel := context.WithTimeout(ctx, grpcStopTimeout)
						defer cancel()
						err := grpcServer.Shutdown(ctx)
						if err != nil {
							logger.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
							return err
						}
						return nil
					},
				},
			)
		},
	),
)
