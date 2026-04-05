package storage

import "github.com/redis/go-redis/v9"

type taskRedisStorage struct {
	client *redis.Client
}

func NewTaskRedisStorage(client *redis.Client) *taskRedisStorage {
	return &taskRedisStorage{
		client: client,
	}
}
