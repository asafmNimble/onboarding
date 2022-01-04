package redis

/*
import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"onboarding/common/data/entities"
	"strconv"
	"time"
)

type RedisTraffic struct {
	GenericRedisBackend
	client   redis.UniversalClient
	zsetName string
}

func NewRedisTraffic() *RedisTraffic {
	rand.Seed(time.Now().UnixNano())
	return &RedisTraffic{client: GetRedis(), zsetName: "awesome_traffic"}
}

// this function increases traffic for an account and project, it's logic:
// put key in zset with score of minutes and seconds if it doesnt exist
// if it's not ready to be processed yet, put it back with new time
// flush will also update the time in the zset
// also use incr for a key with the same name to increase traffic
func (t *RedisTraffic) IncreaseTraffic(ctx context.Context, e *entities.Traffic, requestID string) error {

	// this will change when we add more service rate units
	KeyFormat := e.StringKey()
	AccountProjectTrafficKey := &redis.Z{Score: float64(time.Now().Unix()), Member: KeyFormat}

	// set a key with the request id to make sure it only get's processed once
	resultBool := t.client.SetNX(ctx, requestID, "1", time.Hour*24)
	wasSet, err := resultBool.Result()
	if err != nil {
		return err
	}

	// if requestId was already set
	if !wasSet {
		return err
	}

	// Increase traffic spent bytes for the account and project
	err = t.IncrementTraffic(ctx, KeyFormat, e.SpentBytes)
	if err != nil {
		return err
	}

	// Create a member in the sorted set with a score of the current time, only if it doesn't exist already
	result := t.client.ZAddNX(ctx, t.zsetName, AccountProjectTrafficKey)
	err = result.Err()
	if err != nil {
		return err
	}

	return nil
}

// this function waits until the zset's length is bigger than minLength, sleeping for duration in between
func (t *RedisTraffic) waitForZSetLength(ctx context.Context, minLength int64, duration time.Duration) int64 {
	// get length of sorted set, if smaller than 10, sleep.
	length, err := t.getZsetLength(ctx)
	if err != nil {
		return 1
	}
	for length < minLength {
		time.Sleep(duration)
		length, err = t.getZsetLength(ctx)
		if err != nil {
			return 1
		}
	}
	return length
}

func (t *RedisTraffic) getZsetLength(ctx context.Context) (int64, error) {
	resultInt := t.client.ZCount(ctx, t.zsetName, "-inf", "+inf")
	length, err := resultInt.Result()
	if err != nil {
		fmt.Println("error getting length of sorted set when flushing, continuing.")
		return 0, errors.New("error getting length of sorted set")
	}
	return length, nil
}

func (t *RedisTraffic) GetAccProjScore(ctx context.Context) (accountProjectKey string, timestamp int64, err error) {
	// Sleep until length of sorted set with keys is at least 1
	length := t.waitForZSetLength(ctx, 1, time.Second)
	var n int64
	if length < 10 {
		n = length
	} else {
		n = 9
	}

	// random num between 1 and 9 to get a random lowest score key
	randomNum := rand.Intn(int(n))

	// get key with random Num index from the sorted set with it's associated score
	result := t.client.ZRangeWithScores(ctx, t.zsetName, int64(randomNum), int64(randomNum))
	accountProjectKeySlice, err := result.Result()
	if err != nil {
		return "", 0, err
	}
	// Verify we got only one key from the set
	sliceLen := len(accountProjectKeySlice)
	if sliceLen != 1 {
		return "",
			0,
			fmt.Errorf("expected one account:project key from zset, got: %v", sliceLen)
	}
	accountProjectKey, timestamp = accountProjectKeySlice[0].Member.(string), int64(accountProjectKeySlice[0].Score)
	return accountProjectKey, timestamp, err
}

func (t *RedisTraffic) GetAndResetAccProjTraffic(ctx context.Context,
	accountProjectKey string) (spentBytes int64, err error) {
	resultString := t.client.GetSet(ctx, accountProjectKey, "0")
	spentBytesString, err := resultString.Result()
	if err != nil {
		return 0, err
	}
	// parse spentBytes string to int64
	spentBytes, err = strconv.ParseInt(spentBytesString, 10, 64)
	if err != nil {
		// increment the same back and return with error
		err2 := t.IncrementTraffic(ctx, accountProjectKey, spentBytes)
		if err2 != nil {
			return 0, errors.New(err.Error() + err2.Error())
		}
		return 0, err
	}
	return spentBytes, nil
}

// Perform a redis INCRBY
func (t *RedisTraffic) IncrementTraffic(ctx context.Context, accountProjectKey string, spentBytes int64) error {
	result := t.client.IncrBy(ctx, accountProjectKey, spentBytes)
	return result.Err()
}

func (t *RedisTraffic) ResetTrafficTimeScore(ctx context.Context, accountProjectKey string) error {
	accountProjectTrafficKey := &redis.Z{Score: float64(time.Now().Unix()), Member: accountProjectKey}
	// Create a member in the sorted set with a score of the current time, only if it doesn't exist already
	result := t.client.ZAdd(ctx, t.zsetName, accountProjectTrafficKey)
	err := result.Err()
	if err != nil {
		return err
	}
	return nil
}
*/
