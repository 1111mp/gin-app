package router

import (
	"bytes"
	"io"
	"time"

	"github.com/1111mp/gin-app/config"
	_ "github.com/1111mp/gin-app/docs"
	"github.com/1111mp/gin-app/internal/middleware"
	"github.com/1111mp/gin-app/internal/observability/tracing"
	"github.com/1111mp/gin-app/pkg/logger"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	timeout "github.com/vearne/gin-timeout"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewRouter -.
// Swagger spec:
// @title       Gin APP API
// @description This is a sample server Petstore server.
// @version 		1.0
// @host 				localhost:4000

// @securityDefinitions.apikey APIAuth
// @in												cookie
// @name											app_cookie_name

// @securityDefinitions.apikey OpenAPIAuth
// @in												header
// @name											PRIVATE-TOKEN
func NewRouter(
	cfg config.Config,
	logger logger.Logger,
	filter tracing.TraceFilter,
) *gin.Engine {
	app := gin.New()

	// apply middlewares
	//  OpenTelemetry tracing + metrics
	app.Use(
		otelgin.Middleware(
			cfg.App().Name,
			otelgin.WithMeterProvider(otel.GetMeterProvider()),
			otelgin.WithGinFilter(func(ctx *gin.Context) bool {
				return filter.ShouldTrace(ctx.Request.URL.Path)
			}),
		),
	)
	// request ID
	app.Use(requestid.New())
	// CORS
	app.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Request-ID", "PRIVATE-TOKEN"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:8080"
		},
		MaxAge: 1 * time.Hour,
	}))
	// logging
	app.Use(ginzap.GinzapWithConfig(logger.Logger(), &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Skipper: func(ctx *gin.Context) bool {
			return filter.ShouldTrace(ctx.Request.URL.Path)
		},
		Context: func(ctx *gin.Context) []zapcore.Field {
			var fields []zapcore.Field
			// log request ID
			if rid := requestid.Get(ctx); rid != "" {
				fields = append(fields, zap.String("request_id", rid))
			}

			// log trace and span ID
			span := trace.SpanFromContext(ctx.Request.Context())
			sc := span.SpanContext()
			if sc.IsValid() {
				fields = append(fields, zap.String("trace_id", sc.TraceID().String()))
				fields = append(fields, zap.String("span_id", sc.SpanID().String()))
			}

			if cfg.App().IsDev() {
				// log request body
				bodyBytes, err := ctx.GetRawData()
				if err == nil {
					fields = append(fields, zap.String("body", string(bodyBytes)))
					ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}

			return fields
		},
	}))
	// recovery
	app.Use(ginzap.RecoveryWithZap(logger.Logger(), true))
	// sentry
	app.Use(sentrygin.New(sentrygin.Options{
		Repanic:         true,
		WaitForDelivery: false,
		Timeout:         2 * time.Second,
	}))
	// timeout
	app.Use(timeout.Timeout(timeout.WithTimeout(6 * time.Second)))
	// custom error handler
	app.Use(middleware.ErrorHandler(logger))

	// Swagger
	if cfg.Swagger().Enabled {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return app
}
