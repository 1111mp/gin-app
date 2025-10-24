package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DeRuina/timberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Interface -.
type Interface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(dir string, level string) *Logger {
	var l zapcore.Level

	switch strings.ToLower(level) {
	case "error":
		l = zap.ErrorLevel
	case "warn":
		l = zap.WarnLevel
	case "info":
		l = zap.InfoLevel
	case "debug":
		l = zap.DebugLevel
	default:
		l = zap.InfoLevel
	}

	var core zapcore.Core

	// console
	consoleCfg := zap.NewDevelopmentEncoderConfig()
	consoleCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)

	if dir != "" {
		timberLogger := &timberjack.Logger{
			Filename:           filepath.Join(dir, "gin-app.log"),
			MaxSize:            200,
			MaxBackups:         3,
			MaxAge:             14,
			Compression:        "gzip",
			LocalTime:          true,
			RotationInterval:   24 * time.Hour,
			RotateAt:           []string{"00:00", "12:00"},
			BackupTimeFormat:   "2006-01-02-15-04-05",
			AppendTimeAfterExt: true,
			// RotateAtMinutes:    []int{0, 15, 30, 45},
		}

		// file
		fileCfg := zap.NewProductionEncoderConfig()
		fileCfg.EncodeTime = zapcore.RFC3339TimeEncoder
		fileEncoder := zapcore.NewJSONEncoder(fileCfg)

		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), l),
			zapcore.NewCore(fileEncoder, zapcore.AddSync(timberLogger), l),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), l),
		)
	}

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Logger{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

func (l *Logger) Logger() *zap.Logger {
	return l.logger
}

// Debug -.
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Info -.
func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Warn -.
func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

// Error -.
func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Fatal -.
func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

// Debugf -.
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

// Infof -.
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

// Warnf -.
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

// Errorf -.
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

// Fatalf -.
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}
