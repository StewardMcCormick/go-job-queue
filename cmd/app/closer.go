package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type dependency struct {
	name  string
	close func(ctx context.Context) error
}

type closer struct {
	log        *zap.Logger
	closeFuncs []dependency
}

func NewCloser() *closer {
	return &closer{
		closeFuncs: make([]dependency, 0),
	}
}

func (c *closer) Add(name string, fn func(ctx context.Context) error) {
	c.closeFuncs = append(c.closeFuncs, dependency{name: name, close: fn})
}

func (c *closer) Close(ctx context.Context) error {
	globalStart := time.Now()
	c.log.Info("[SHUTDOWN] Start shutting down...")
	errs := make([]error, 0)
	for i := len(c.closeFuncs) - 1; i >= 0; i-- {
		dep := c.closeFuncs[i]

		c.log.Info(fmt.Sprintf("[SHUTDOWN] %s closing...", dep.name))
		start := time.Now()
		if err := dep.close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("[SHUTDOWN] Cannot close %s: %w", dep.name, err))
			c.log.Error(fmt.Sprintf("[SHUTDOWN] Cannot close %s: %v", dep.name, err))
			continue
		}
		c.log.Info(fmt.Sprintf("[SHUTDOWN] %s closed, duration: %d ms", dep.name, time.Since(start).Milliseconds()))
	}
	c.log.Info(fmt.Sprintf("[SHUTDOWN] Shutdown completed. Total duration: %d ms",
		time.Since(globalStart).Milliseconds()),
	)

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}
