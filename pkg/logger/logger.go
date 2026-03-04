package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DeRuina/timberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger -.
type Logger interface {
	Logger() *zap.Logger

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

	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

// loggerImpl -.
type loggerImpl struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

var _ Logger = (*loggerImpl)(nil)

// New -.
func New(dir string, level string) Logger {
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

	cores := []zapcore.Core{
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), l),
	}

	if dir != "" {
		// file
		fileCfg := zap.NewProductionEncoderConfig()
		fileCfg.EncodeTime = zapcore.RFC3339TimeEncoder
		fileEncoder := zapcore.NewJSONEncoder(fileCfg)

		appLogWriter := getLogWriter(filepath.Join(dir, "app.log"))
		appCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(appLogWriter), l)
		cores = append(cores, appCore)

		errLogWriter := getLogWriter(filepath.Join(dir, "app-error.log"))
		errCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(errLogWriter), zap.ErrorLevel)
		cores = append(cores, errCore)
	}

	core = zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &loggerImpl{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

func getLogWriter(filename string) io.Writer {
	return &timberjack.Logger{
		Filename:           filename,
		MaxSize:            200, // MB
		MaxBackups:         3,
		MaxAge:             14, // Days
		Compression:        "gzip",
		LocalTime:          true,
		RotationInterval:   24 * time.Hour,
		RotateAt:           []string{"00:00", "12:00"},
		BackupTimeFormat:   "2006-01-02-15-04-05",
		AppendTimeAfterExt: true,
		// RotateAtMinutes:    []int{0, 15, 30, 45},
	}
}

func (l *loggerImpl) Logger() *zap.Logger {
	return l.logger
}

// Debug -.
func (l *loggerImpl) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Info -.
func (l *loggerImpl) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Warn -.
func (l *loggerImpl) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

// Error -.
func (l *loggerImpl) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Fatal -.
func (l *loggerImpl) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

// Debugf -.
func (l *loggerImpl) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

// Infof -.
func (l *loggerImpl) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

// Warnf -.
func (l *loggerImpl) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

// Errorf -.
func (l *loggerImpl) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

// Fatalf -.
func (l *loggerImpl) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

// Infow -.
func (l *loggerImpl) Infow(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Errorw -.
func (l *loggerImpl) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}
