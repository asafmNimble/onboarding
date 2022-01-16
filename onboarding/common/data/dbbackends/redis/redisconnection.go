package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strings"
	"sync"
)

var once sync.Once
var singletonClient redis.UniversalClient

func remoteRedisConnect() (redis.UniversalClient, error) {
	hostSlice := strings.Split("172.38.0.63:6379", ",")

	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    hostSlice,
		Username: "",
		Password: ""})

	return redisClient, nil
}

func GetRedis(clients ...redis.UniversalClient) redis.UniversalClient {

	once.Do(func() {
		var client redis.UniversalClient
		var err error
		// If passed a client, update the singleton to hold that client.
		if len(clients) > 0 {
			client = clients[0]
		} else {
			// otherwise, get remote redis
			client, err = remoteRedisConnect()
			if err != nil {
				fmt.Println("error connecting to redis", zap.String("error", err.Error()))
			}
		}
		singletonClient = client
	})
	return singletonClient
}
