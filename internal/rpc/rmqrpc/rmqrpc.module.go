package rmqrpc

import (
	"context"
	"fmt"

	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/pkg/logger"
	rmqRPCClient "github.com/1111mp/gin-app/pkg/rabbitmq/rmq_rpc/client"
	rmqRPCServer "github.com/1111mp/gin-app/pkg/rabbitmq/rmq_rpc/server"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"rmqrpc",

	fx.Provide(
		// rmqrpc controller
		NewRMQRPCControllter,
		// rabbitmq rpc server
		func(
			cfg config.Config,
			logger logger.Logger,
		) (*rmqRPCServer.Server, error) {
			rmqServer, err := rmqRPCServer.New(cfg.RMQ().URL, cfg.RMQ().ServerExchange, logger)
			if err != nil {
				logger.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
				return nil, err
			}

			logger.Info("app - Run - rmqServer - initialized")
			return rmqServer, nil
		},
		// rabbitmq rpc client
		fx.Annotate(
			func(cfg config.Config, logger logger.Logger) (*rmqRPCClient.Client, error) {
				rmqClient, err := rmqRPCClient.New(
					cfg.RMQ().URL,
					cfg.RMQ().ServerExchange,
					cfg.RMQ().ClientExchange,
				)
				if err != nil {
					logger.Fatal(fmt.Errorf("app - Run - rmqClient - client.New: %w", err))
					return nil, err
				}

				logger.Info("app - Run - rmqClient - initialized")
				return rmqClient, nil
			},
			fx.OnStop(func(logger logger.Logger, rmqClient *rmqRPCClient.Client) error {
				logger.Info("app - Run - rmqClient - shutting down")

				err := rmqClient.Shutdown()
				if err != nil {
					logger.Error(fmt.Errorf("app - Run - rmqClient - client.Shutdown: %w", err))
					return err
				}

				logger.Info("app - Run - rmqClient - shutdown completed")
				return nil
			}),
		),
	),
	// register router for rabbitmq
	fx.Invoke(
		func(rmqServer *rmqRPCServer.Server, c RMQRPCControllter) {
			rmqServer.RegisterRouter("v1.get_post_by_id", rmqRPCServer.CallHandlerTyped(c.GetPostById))
		},
	),

	fx.Invoke(func(
		lc fx.Lifecycle,
		sd fx.Shutdowner,
		logger logger.Logger,
		rmqServer *rmqRPCServer.Server,
	) {
		lc.Append(
			fx.Hook{
				OnStart: func(ctx context.Context) error {
					rmqServer.Start()

					go func() {
						for {
							select {
							case <-ctx.Done():
								return
							case err, ok := <-rmqServer.Notify():
								if !ok {
									return
								}

								logger.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
								// ! It's better to just shut down the app and let Kubernetes restart it.
								_ = sd.Shutdown()
								return
							}
						}
					}()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					if err := rmqServer.Shutdown(); err != nil {
						logger.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
						return err
					}
					return nil
				},
			},
		)
	}),
)
