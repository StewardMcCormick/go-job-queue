package appctx

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const LoggerKey contextKey = "app_ctx_logger"

func WithLogger(parent context.Context, log *zap.Logger) context.Context {
	return context.WithValue(parent, LoggerKey, log)
}

func GetLogger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.L()
	}

	log, ok := ctx.Value(LoggerKey).(*zap.Logger)
	if !ok {
		return zap.L()
	}

	return log
}
