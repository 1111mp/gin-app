package router

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/1111mp/gin-app/config"
	_ "github.com/1111mp/gin-app/docs"
	api "github.com/1111mp/gin-app/internal/api/v1"
	"github.com/1111mp/gin-app/internal/service"
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
// @title       Go Clean Template API
// @description This is a sample server Petstore server.
// @version 		1.0
// @host 				localhost:8080
// @BasePath 		/api/v1
func NewRouter(app *gin.Engine, cfg *config.Config, l *logger.Logger) {
	// middleware
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:8080"
		},
		MaxAge: 12 * time.Hour,
	}))

	app.Use(ginzap.GinzapWithConfig(l.Logger(), &ginzap.Config{
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
	app.Use(ginzap.RecoveryWithZap(l.Logger(), true))

	app.Use(timeout.Timeout(timeout.WithTimeout(3 * time.Second)))

	// Swagger
	if cfg.Swagger.Enabled {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// K8s probe
	app.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s := service.NewServiceGroup(l)
	a := api.NewApiGroup(s)
	r := NewRouterGroup(a)

	// Routes
	apiV1Group := app.Group("/api/v1")
	{
		r.UserRouter.RegisterRoutes(apiV1Group)
	}
}
