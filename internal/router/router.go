package router

import (
	"bytes"
	"io"
	"net/http"
	"time"

	api "github.com/1111mp/gin-app/internal/api/v1"
	"github.com/1111mp/gin-app/internal/service"
	"github.com/1111mp/gin-app/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	timeout "github.com/vearne/gin-timeout"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewRouter -.
func NewRouter(router *gin.Engine, l *logger.Logger) {
	// middleware
	router.Use(requestid.New())
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:8080"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.Use(ginzap.GinzapWithConfig(l.Logger(), &ginzap.Config{
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
	router.Use(ginzap.RecoveryWithZap(l.Logger(), true))

	router.Use(timeout.Timeout(timeout.WithTimeout(3 * time.Second)))

	// K8s probe
	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s := service.NewServiceGroup(l)
	a := api.NewApiGroup(s)
	r := NewRouterGroup(a)

	// Routes
	apiV1Group := router.Group("/api/v1")
	{
		r.UserRouter.RegisterRoutes(apiV1Group)
	}
}
