package ctxlog

import (
	"context"

	"go.uber.org/zap"
)

type ctxLoggerKey struct{}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, logger)
}

func Logger(ctx context.Context) *zap.Logger {
	return ctx.Value(ctxLoggerKey{}).(*zap.Logger)
}
