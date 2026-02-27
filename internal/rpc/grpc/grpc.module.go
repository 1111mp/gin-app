package grpc

import (
	"context"
	"fmt"

	"github.com/1111mp/gin-app/config"
	v1 "github.com/1111mp/gin-app/docs/proto/v1"
	grpcclient "github.com/1111mp/gin-app/pkg/grpc/client"
	"github.com/1111mp/gin-app/pkg/grpcserver"
	"github.com/1111mp/gin-app/pkg/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc/reflection"
)

var Module = fx.Module(
	"grpc",

	fx.Provide(
		// grpc controller
		NewGRPCControllter,
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
			grpcController GRPCControllter,
		) {
			{
				v1.RegisterPostServer(grpcServer.App, grpcController)
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
							for {
								select {
								case <-ctx.Done():
									return
								case err, ok := <-grpcServer.Notify():
									if !ok {
										return
									}

									logger.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
									// ! It's better to just shut down the app and let Kubernetes restart it.
									_ = sd.Shutdown()
									return
								}
							}
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						err := grpcServer.Shutdown()
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
