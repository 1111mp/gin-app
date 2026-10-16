package rmqrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/pkg/logger"
	rmqRPCClient "github.com/1111mp/gin-app/pkg/rabbitmq/rmq_rpc/client"
	rmqRPCServer "github.com/1111mp/gin-app/pkg/rabbitmq/rmq_rpc/server"
	"go.uber.org/fx"
)

const (
	serverStopTimeout = 5 * time.Second
	clientStopTimeout = 3 * time.Second
)

var Module = fx.Module(
	"rmqrpc",

	fx.Provide(
		// rmqrpc router
		NewRMQRPCRouter,
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
			fx.OnStop(func(ctx context.Context, logger logger.Logger, rmqClient *rmqRPCClient.Client) error {
				logger.Info("app - Run - rmqClient - shutting down")

				ctx, cancel := context.WithTimeout(ctx, clientStopTimeout)
				defer cancel()
				err := rmqClient.Shutdown(ctx)
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
		func(rmqServer *rmqRPCServer.Server, rmqrpcRouter RMQRPCRouter) {
			rmqServer.RegisterRouter("v1.get_post_by_id", rmqRPCServer.CallHandlerTyped(rmqrpcRouter.GetPostById))
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
						if err, ok := <-rmqServer.Notify(); ok {
							logger.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
							if err := sd.Shutdown(fx.ExitCode(1)); err != nil {
								logger.Error(err)
							}
						}
					}()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					ctx, cancel := context.WithTimeout(ctx, serverStopTimeout)
					defer cancel()
					if err := rmqServer.Shutdown(ctx); err != nil {
						logger.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
						return err
					}
					return nil
				},
			},
		)
	}),
)
