package appctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetLogger(t *testing.T) {
	t.Run("returns logger set via WithLogger", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()
		defer logger.Sync()

		ctx := WithLogger(context.Background(), logger)
		result := GetLogger(ctx)

		assert.Equal(t, logger, result)
	})

	t.Run("returns global logger when context has no logger", func(t *testing.T) {
		ctx := context.Background()
		result := GetLogger(ctx)

		assert.Equal(t, zap.L(), result)
	})

	t.Run("returns global logger when context has wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), LoggerKey, "not_a_logger")
		result := GetLogger(ctx)

		assert.Equal(t, zap.L(), result)
	})

	t.Run("returns global logger when context is nil", func(t *testing.T) {
		var ctx context.Context
		result := GetLogger(ctx)

		assert.Equal(t, zap.L(), result)
	})

	t.Run("handles nil logger in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), LoggerKey, nil)
		result := GetLogger(ctx)

		// Type assertion fails, returns global logger
		assert.Equal(t, zap.L(), result)
	})
}
