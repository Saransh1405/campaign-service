package logger

import (
	"context"
	"users-service/constants"

	"go.uber.org/zap"
)

var log Logger

type Logger interface {
	Info(string, ...zap.Field)
	Warn(string, ...zap.Field)
	Error(string, ...zap.Field)
	Debug(string, ...zap.Field)
	Fatal(string, ...zap.Field)
	Panic(string, ...zap.Field)
	With(args ...zap.Field) Logger
	Sync()
}

func GetLogger(ctx context.Context) Logger {
	return log.With(zap.String(constants.RequestId, ctx.Value(constants.RequestId).(string)))
}

func GetLoggerWithoutContext() Logger {
	return log
}
