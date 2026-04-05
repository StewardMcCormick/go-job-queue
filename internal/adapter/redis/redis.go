package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host        string        `env:"REDIS_HOST" required:"true"`
	Port        string        `env:"REDIS_PORT" required:"true"`
	PoolSize    int           `yaml:"pool_size" env-default:"10"`
	PoolTimeout time.Duration `yaml:"pool_timeout" env-default:"5s"`
	Password    string        `env:"REDIS_PASSWORD" required:"true"`
}

func NewConnection(cfg Config, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  cfg.PoolTimeout,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		DB:           db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err.Err() != nil {
		closingErr := client.Close()
		if closingErr != nil {
			return nil, fmt.Errorf("redis get connection and closing error: %w, %w", err.Err(), closingErr)
		}
		return nil, fmt.Errorf("redis get connection error: %w", err.Err())
	}

	return client, nil
}
