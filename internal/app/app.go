// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/1111mp/gin-app/config"
	"github.com/1111mp/gin-app/internal/modules"
	"github.com/1111mp/gin-app/internal/router"
	api_router "github.com/1111mp/gin-app/internal/router/api"
	openapi_router "github.com/1111mp/gin-app/internal/router/open-api"
	"github.com/1111mp/gin-app/internal/rpc"
	"github.com/1111mp/gin-app/pkg/httpserver"
	"github.com/1111mp/gin-app/pkg/jwt"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/1111mp/gin-app/pkg/oauth2/github"
	"github.com/1111mp/gin-app/pkg/oauth2/google"
	"github.com/1111mp/gin-app/pkg/postgres"
	"github.com/1111mp/gin-app/pkg/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Run creates objects via constructors.
func Run(cfg config.Config) { //nolint: gocyclo,cyclop,funlen,gocritic,nolintlint
	fx.New(
		fx.Supply(
			fx.Annotate(
				cfg,
				fx.As(new(config.Config)),
			),
		),

		fx.Provide(
			// logger
			func(cfg config.Config) logger.Logger {
				return logger.New(cfg.Log().Dir, cfg.Log().Level)
			},
			// postgres
			fx.Annotate(
				func(cfg config.Config, logger logger.Logger) (*postgres.Postgres, error) {
					pg, err := postgres.New(cfg.PG().URL, postgres.MaxPoolSize(cfg.PG().PoolMax))
					if err != nil {
						logger.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
						return nil, err
					}

					logger.Info("app - Run - postgres connected")
					return pg, nil
				},
				fx.OnStop(func(logger logger.Logger, pg *postgres.Postgres) error {
					logger.Info("app - Run - postgres closed")
					pg.Close()
					return nil
				}),
			),
			// redis
			fx.Annotate(
				func(cfg config.Config, logger logger.Logger) (*redis.Redis, error) {
					rdb, err := redis.New(cfg.Redis().URL, redis.MaxPoolSize(cfg.Redis().PoolMax))
					if err != nil {
						logger.Fatal(fmt.Errorf("app - Run - redis.New: %w", err))
						return nil, err
					}

					logger.Info("app - Run - redis connected")
					return rdb, nil
				},
				fx.OnStop(func(logger logger.Logger, rdb *redis.Redis) error {
					logger.Info("app - Run - redis closed")

					rdb.Close()
					return nil
				}),
			),
			// jwt
			func() jwt.JWT {
				return jwt.NewJWT(jwt.Issuer(cfg.App().Name), jwt.Secret(cfg.JWT().SECRET))
			},
			// github oauth2
			func() github.Client {
				return github.Setup(
					github.ClientID(cfg.Github().ClientID),
					github.ClientSecret(cfg.Github().ClientSecret),
					github.RedirectURL(cfg.Github().RedirectURL),
				)
			},
			// google oauth2
			func() google.Client {
				return google.Setup(
					google.ClientID(cfg.Google().ClientID),
					google.ClientSecret(cfg.Google().ClientSecret),
					google.RedirectURL(cfg.Google().RedirectURL),
				)
			},
			// gin
			router.NewRouter,
			// api router
			api_router.NewAPIRouter,
			// open-api router
			openapi_router.NewOpenAPIRouter,
			// http server
			func(cfg config.Config, handler *gin.Engine) httpserver.Server {
				return httpserver.New(handler, httpserver.Port(cfg.HTTP().Port))
			},
		),
		// rpc
		rpc.Module,
		// api
		modules.APIModule,
		// start http server
		fx.Invoke(startHTTPServer),
	).Run()
}

func startHTTPServer(
	lc fx.Lifecycle,
	sd fx.Shutdowner,
	logger logger.Logger,
	httpServer httpserver.Server,
) {
	lc.Append(
		fx.Hook{
			// start
			OnStart: func(ctx context.Context) error {
				logger.Info("http server - Starting HTTP Server...")

				go func() {
					if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
						logger.Errorf("http server - Start HTTP Server error: %v", err)
						sd.Shutdown()
					}
				}()

				logger.Info("http server - Listening on ", httpServer.GetAddress())
				return nil
			},
			// stop
			OnStop: func(ctx context.Context) error {
				logger.Info("http server - Stopping HTTP Server...")

				err := httpServer.Shutdown()
				if err != nil {
					return err
				}

				logger.Info("http server - Server - Shutting down")
				return nil
			},
		},
	)
}
