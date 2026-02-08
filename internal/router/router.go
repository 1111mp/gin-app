package router

import (
	"bytes"
	"io"
	"time"

	"github.com/1111mp/gin-app/config"
	_ "github.com/1111mp/gin-app/docs"
	"github.com/1111mp/gin-app/internal/middleware"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	timeout "github.com/vearne/gin-timeout"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewRouter -.
// Swagger spec:
// @title       Gin APP API
// @description This is a sample server Petstore server.
// @version 		1.0
// @host 				localhost:8080

// @securityDefinitions.apikey APIAuth
// @in												cookie
// @name											app_cookie_name

// @securityDefinitions.apikey OpenAPIAuth
// @in												header
// @name											PRIVATE-TOKEN
func NewRouter(
	cfg config.ConfigInterface,
	logger logger.Interface,
) *gin.Engine {
	app := gin.Default()

	// apply middlewares
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Request-ID", "PRIVATE-TOKEN"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:8080"
		},
		MaxAge: 1 * time.Hour,
	}))

	app.Use(ginzap.GinzapWithConfig(logger.Logger(), &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: func(ctx *gin.Context) []zapcore.Field {
			var fields []zapcore.Field
			// log request ID
			if rid := requestid.Get(ctx); rid != "" {
				fields = append(fields, zap.String("request_id", rid))
			}

			// log trace and span ID
			if trace.SpanFromContext(ctx.Request.Context()).SpanContext().IsValid() {
				fields = append(fields, zap.String("trace_id", trace.SpanFromContext(ctx.Request.Context()).SpanContext().TraceID().String()))
				fields = append(fields, zap.String("span_id", trace.SpanFromContext(ctx.Request.Context()).SpanContext().SpanID().String()))
			}

			// log request body
			var body []byte
			var buf bytes.Buffer
			tee := io.TeeReader(ctx.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			ctx.Request.Body = io.NopCloser(&buf)
			fields = append(fields, zap.String("body", string(body)))

			return fields
		},
	}))
	app.Use(ginzap.RecoveryWithZap(logger.Logger(), true))
	app.Use(timeout.Timeout(timeout.WithTimeout(6 * time.Second)))
	app.Use(middleware.ErrorHandler(logger))

	// Swagger
	if cfg.Swagger().Enabled {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return app
}
