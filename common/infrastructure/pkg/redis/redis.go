package redis

import (
	"context"
	"example-go/common/global"
	"example-go/common/infrastructure/config"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewRedis(config *config.ServerConfig) *redis.Client {
	c := config.DecimalCricketRedis
	client := redis.NewClient(&redis.Options{
		Addr:         c.Address,
		DB:           0,
		PoolSize:     config.DecimalCricketRedis.PoolSize,
		DialTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		MinIdleConns: config.DecimalCricketRedis.MinIdleConns,
	})

	client.AddHook(DecimalCricketHook{})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	global.Redis = client
	return client
}

func NewRedisReadOnly(config *config.ServerConfig) *redis.Client {
	c := config.DecimalCricketRedisReadOnly

	client := redis.NewClient(&redis.Options{
		Addr:         c.Address,
		DB:           0,
		PoolSize:     config.DecimalCricketRedisReadOnly.PoolSize,
		DialTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		MinIdleConns: config.DecimalCricketRedisReadOnly.MinIdleConns,
	})

	client.AddHook(DecimalCricketReadOnlyHook{})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	global.RedisReadOnly = client
	return client
}
