package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.Logger

func Init(env string) {
	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
	}
	l, _ := cfg.Build()
	defaultLogger = l
}

func L() *zap.Logger { return defaultLogger }

func WithContext(ctx context.Context) *zap.Logger {
	// trace_id 提取可后续接入 tracer 包
	return defaultLogger
}
