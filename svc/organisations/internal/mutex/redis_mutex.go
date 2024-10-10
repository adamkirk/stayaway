package mutex

import (
	"context"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	lock *redislock.Lock
}

func (rl *RedisLock) Release() error {
	err := rl.lock.Release(context.TODO())

	if err == nil || err == redislock.ErrLockNotHeld {
		return nil
	}

	return err
}

type RedisConnector interface {
	Client() (*redis.Client, error)
}

type RedisMutex struct {
	client *redislock.Client
	connector RedisConnector
}

func (rm *RedisMutex) getClient() (*redislock.Client, error) {
	if rm.client != nil {
		return rm.client, nil
	}

	redisClient, err := rm.connector.Client()

	if err != nil {
		return nil, err
	}

	rm.client = redislock.New(redisClient)

	return rm.client, nil
}

func (rm *RedisMutex) prefixKey(key string) string {
	return fmt.Sprintf("%s:%s", "organisations_mutex", key)
}

// ClaimWithBackOff claims a a lock for the given key and retries it 3 times
// with a 100 ms interval between. This seems a sensible default for most use
// cases in the app.
func (rm *RedisMutex) ClaimWithBackOff(key string, ttl time.Duration) (DistributedMutex, error) {
	client, err := rm.getClient()

	if err != nil {
		return nil, err
	}

	lockKey := rm.prefixKey(key)

	l, err := client.Obtain(context.TODO(), lockKey, ttl, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(100 * time.Millisecond), 3),
	})

	if err != nil {
		return nil, ErrLockNotClaimed{
			Key: key,
			Err: err,
		}
	}

	return &RedisLock{
		lock: l,
	}, nil
}

func NewRedisMutex(connector RedisConnector) *RedisMutex {
	return &RedisMutex{
		connector: connector,
	}
}