package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	//level := config.AppConfig.LogLevel
	log, err := zap.NewProduction()
	fmt.Println(err)
	return log
}

func InitLogger(level string) (*zap.Logger, error) {
	logConfig := zap.Config{

		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
		},
	}

	switch level {
	case "debug":
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		break
	case "warn":
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		break
	case "error":
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		break
	case "fatal":
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
		break
	}

	return logConfig.Build()
}
