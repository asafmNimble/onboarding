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

func (s *RedisGuessersCounter) CreateGuessersCounter(gc *entities.GuesserCounter) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	gID := strconv.Itoa(int(gc.GuesserID))
	gCount := strconv.Itoa(int(gc.Counter))
	_, err := s.client.SetNX(ctx, gID, gCount, 5*time.Minute).Result()
	return err
}

func (s *RedisGuessersCounter) IncreaseGuesserCounter(guesserID int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	gID := strconv.Itoa(int(guesserID))
	counter, err := s.client.Incr(ctx, gID).Result()
	if err != nil {
		return 0, err
	}
	s.client.Expire(context.Background(), gID, 5*time.Minute)
	return counter, nil
}

func (s *RedisGuessersCounter) GetGuesserCounter(guesserID int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	gID := strconv.Itoa(int(guesserID))
	valStr, err := s.client.Get(ctx, gID).Result()
	if err != nil {
		return -1, err
	}
	s.client.Expire(context.Background(), gID, 5*time.Minute)
	val, _ := strconv.Atoi(valStr)
	return int64(val), nil
}
