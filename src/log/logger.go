package log

import (
	"fmt"

	"go.uber.org/zap"
)

var instance Logger

type Logger interface {
	Info(msg string, ctxs ...LoggerContext[any])
	Debug(msg string, ctxs ...LoggerContext[any])
	Error(err error, ctxs ...LoggerContext[any])
}

// NewLogger return singleton Logger instance
func NewLogger() Logger {
	if instance == nil {
		instance = NewZapLogger()
	}
	return instance
}

type LoggerContext[T any] struct {
	key string
	val T
}
type StringLoggerContext = LoggerContext[string]

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() *ZapLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return &ZapLogger{logger}
}

func (z *ZapLogger) Info(msg string, ctxs ...LoggerContext[any]) {
	fields := z.ctxsToFields(ctxs...)
	z.logger.Info(msg, fields...)
}

func (z *ZapLogger) Debug(msg string, ctxs ...LoggerContext[any]) {
	fields := z.ctxsToFields(ctxs...)
	z.logger.Debug(msg, fields...)
}

func (z *ZapLogger) Error(err error, ctxs ...LoggerContext[any]) {
	fields := z.ctxsToFields(ctxs...)
	z.logger.Error(err.Error(), fields...)
}

func (z *ZapLogger) ctxsToFields(ctxs ...LoggerContext[any]) []zap.Field {
	toField := func(ctx any) (zap.Field, error) {
		switch c := ctx.(type) {
		case StringLoggerContext:
			return zap.String(c.key, c.val), nil
		default:
			return zap.Any("dummy", "dummy"), fmt.Errorf("undefined logger context: %#+v", c)
		}
	}

	fields := make([]zap.Field, 0)
	for _, ctx := range ctxs {
		if field, err := toField(ctx); err == nil {
			fields = append(fields, field)
		} else {
			z.Error(err)
		}
	}
	return fields
}
