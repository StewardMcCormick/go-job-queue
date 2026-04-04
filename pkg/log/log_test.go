package log

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestNewLogger_DevEnvironment(t *testing.T) {
	cfg := Config{
		Level:   "info",
		Outputs: []string{"stdout"},
	}

	logger, err := NewLogger(cfg, "dev", "test-app", "1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("expected logger to be non-nil")
	}
}

func TestNewLogger_ProductionEnvironment(t *testing.T) {
	cfg := Config{
		Level:   "warn",
		Outputs: []string{"stderr"},
	}

	logger, err := NewLogger(cfg, "production", "test-app", "2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("expected logger to be non-nil")
	}
}

func TestNewLogger_UnknownEnvironment(t *testing.T) {
	cfg := Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	}

	logger, err := NewLogger(cfg, "staging", "test-app", "1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("expected logger to be non-nil")
	}
}

func TestNewLogger_InvalidLevel(t *testing.T) {
	cfg := Config{
		Level:   "invalid-level",
		Outputs: []string{"stdout"},
	}

	logger, err := NewLogger(cfg, "dev", "test-app", "1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("expected logger to be non-nil when level is invalid (should default to info)")
	}
}

func TestNewLogger_EmptyLevel(t *testing.T) {
	cfg := Config{
		Level:   "",
		Outputs: []string{"stdout"},
	}

	logger, err := NewLogger(cfg, "dev", "test-app", "1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("expected logger to be non-nil")
	}
}

func TestNewLogger_MultipleOutputs(t *testing.T) {
	cfg := Config{
		Level:   "debug",
		Outputs: []string{"stdout", "stderr"},
	}

	logger, err := NewLogger(cfg, "dev", "test-app", "1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("expected logger to be non-nil")
	}
}

func TestNewLogger_AllValidLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			cfg := Config{
				Level:   level,
				Outputs: []string{"stdout"},
			}

			logger, err := NewLogger(cfg, "dev", "test-app", "1.0.0")
			if err != nil {
				t.Fatalf("unexpected error for level %s: %v", level, err)
			}
			defer logger.Sync()

			if logger == nil {
				t.Fatal("expected logger to be non-nil")
			}
		})
	}
}

func TestNewLogger_LevelDefaultsToInfoWhenInvalid(t *testing.T) {
	cfg := Config{
		Level:   "not-a-level",
		Outputs: []string{"stdout"},
	}

	logger, err := NewLogger(cfg, "dev", "test-app", "1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer logger.Sync()

	// Verify that the logger was created successfully
	if logger == nil {
		t.Fatal("expected logger to be non-nil")
	}

	// The level should default to info
	expectedLevel := zapcore.InfoLevel
	if logger.Level() != expectedLevel {
		t.Errorf("expected level %v, got %v", expectedLevel, logger.Level())
	}
}
