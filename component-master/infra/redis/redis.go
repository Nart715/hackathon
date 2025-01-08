package redis

import (
	"component-master/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewInitRedisClient(rd *config.RedisConfig) (*RedisClient, error) {
	opt, err := redis.ParseURL(rd.BuildRedisConnectionString())
	if err != nil {
		return nil, err
	}
	opt.PoolSize = rd.MaxIdle
	opt.DialTimeout = rd.DialTimeout
	opt.ReadTimeout = rd.ReadTimeout
	opt.WriteTimeout = rd.WriteTimeout
	opt.Password = rd.Password
	opt.DB = rd.DB

	client := redis.NewClient(opt)
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	slog.Info(fmt.Sprintf("Ping to redis success: %s", pong))
	return &RedisClient{
		client: client,
	}, nil
}

func (rc *RedisClient) Close() {
	if rc.client == nil {
		return
	}
	rc.client.Close()
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (rc *RedisClient) SetNX(key string, value string, expiration time.Duration) error {
	err := rc.client.SetNX(context.Background(), key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisClient) AcquireLock(key string, expiration time.Duration) error {
	err := rc.client.SetNX(context.Background(), key, 1, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := rc.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisClient) Del(ctx context.Context, key string) error {
	err := rc.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
