// Package log TODO
package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZap(lvl, encType string) (*zap.Logger, error) {
	var (
		err    error
		logger *zap.Logger
	)
	config := zap.NewProductionConfig()
	config.Encoding = encType
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.LevelKey = "l"
	config.EncoderConfig.CallerKey = "c"
	config.EncoderConfig.StacktraceKey = "s"

	if err := config.Level.UnmarshalText([]byte(lvl)); err != nil && lvl != "" {
		return nil, err
	}
	logger, err = config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

type ctxKey struct{} // or exported to use outside the package

func CtxWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctxLogger, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return ctxLogger
	}
	return zap.NewNop()
}
