package logger

import (
	"os"
	"users-service/utils"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevelMap = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"error": zap.ErrorLevel,
	"fatal": zap.FatalLevel,
	"panic": zap.PanicLevel,
}

type zapLogger struct {
	Log *zap.Logger
}

func SetupLogging(path string, level string) {
	var logLevel zapcore.Level
	var ok bool
	logLevel, ok = logLevelMap[level]
	if !ok {
		logLevel = zap.InfoLevel
	}
	logFile, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		// rotatelogs.WithMaxAge(time.Duration(24)*time.Hour),
		rotatelogs.WithRotationCount(2),
	)
	if err != nil {
		panic("error setting up logger")
	}

	cfg := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: utils.LevelEncoder,
		TimeKey:     "time",
		EncodeTime:  utils.LogTimeEncoder,
		FunctionKey: "method",
	}

	var cores []zapcore.Core
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(logFile),
		logLevel,
	)
	cores = append(cores, fileCore)
	if os.Getenv("Environment") != "production" {
		consoleErrors := zapcore.Lock(os.Stderr)
		consoleCore := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), consoleErrors, logLevel)
		cores = append(cores, consoleCore)
	}
	tee := zapcore.NewTee(cores...)
	log = &zapLogger{zap.New(tee, zap.AddCaller(), zap.AddCallerSkip(1))}
}

func (l *zapLogger) With(args ...zap.Field) Logger {
	return &zapLogger{
		Log: l.Log.With(args...),
	}
}

func (l *zapLogger) Info(msg string, args ...zap.Field) {
	l.Log.Info(msg, args...)
}

func (l *zapLogger) Warn(msg string, args ...zap.Field) {
	l.Log.Warn(msg, args...)
}

func (l *zapLogger) Debug(msg string, args ...zap.Field) {
	l.Log.Debug(msg, args...)
}

func (l *zapLogger) Error(msg string, args ...zap.Field) {
	l.Log.Error(msg, args...)
}

func (l *zapLogger) Fatal(msg string, args ...zap.Field) {
	l.Log.Fatal(msg, args...)
}

func (l *zapLogger) Panic(msg string, args ...zap.Field) {
	l.Log.Panic(msg, args...)
}

func (l *zapLogger) Sync() {
	l.Log.Sync()
}
