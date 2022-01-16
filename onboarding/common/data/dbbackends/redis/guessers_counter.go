package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisGuessersCounter struct {
	GenericRedisBackend
	client redis.UniversalClient
}

func NewRedisGuessersCounter() *RedisGuessersCounter {
	return &RedisGuessersCounter{client: GetRedis()}
}

func (s *RedisGuessersCounter) CheckIfAlive() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	_, err := s.client.Ping(ctx).Result()
	if err != nil {
		return false
	}
	return true
}


