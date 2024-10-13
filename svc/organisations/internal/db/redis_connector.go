package db

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type RedisConnectorConfig interface {
	RedisHost() string
	RedisPassword() *string
	RedisDb() int
	RedisConnectionRetries() int
}

type RedisConnector struct {
	host     string
	password *string
	db       int
	retries  int
	client   *redis.Client
}

func (rc *RedisConnector) Client() (*redis.Client, error) {
	if rc.client != nil {
		return rc.client, nil
	}
	var client *redis.Client
	var err error

	opts := &redis.Options{
		Addr:     rc.host,
		Password: "",
		DB:       rc.db,
	}

	if rc.password != nil {
		opts.Password = *rc.password
	}

	for i := range rc.retries {
		client = redis.NewClient(opts)

		res := client.Ping(context.TODO())
		err = res.Err()
		if err != nil {
			slog.Warn("failed to connect to redis", "attempt", i+1, "error", err)
			continue
		}
	}

	if err != nil {
		rc.client = client
	}

	return client, err
}

func NewRedisConnector(cfg RedisConnectorConfig) *RedisConnector {
	return &RedisConnector{
		host:     cfg.RedisHost(),
		db:       cfg.RedisDb(),
		password: cfg.RedisPassword(),
		retries:  cfg.RedisConnectionRetries(),
	}
}
