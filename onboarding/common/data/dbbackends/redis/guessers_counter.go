package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"onboarding/common/data/entities"
	"strconv"
	"time"
)

type RedisGuessersCounter struct {
	GenericRedisBackend
	client redis.UniversalClient
}

func NewRedisGuessersCounter() *RedisGuessersCounter {
	return &RedisGuessersCounter{client: GetRedis()}
}

/*
func (s *RedisGuessersCounter) CheckIfAlive() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	_, err := s.client.Ping(ctx).Result()
	if err != nil {
		return false
	}
	return true
}
 */

func (s *RedisGuessersCounter) CreateGuessersCounter(gc *entities.GuesserCounter) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	gID := strconv.Itoa(int(gc.GuesserID))
	gCount := strconv.Itoa(int(gc.Counter))
	_, err := s.client.SetNX(ctx, gID, gCount, 0).Result()
	return err
}

func (s *RedisGuessersCounter) IncreaseGuesserCounter(guesserID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	gID := strconv.Itoa(int(guesserID))
	_, err := s.client.IncrBy(ctx, gID, int64(1)).Result()
	if err != nil {return err}
	return nil
}

func (s *RedisGuessersCounter) GetGuesserCounter(guesserID int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	gID := strconv.Itoa(int(guesserID))
	valStr, err := s.client.Get(ctx, gID).Result()
	if err != nil {
		return -1, err
	}
	val, _ := strconv.Atoi(valStr)
	return int64(val), nil
}