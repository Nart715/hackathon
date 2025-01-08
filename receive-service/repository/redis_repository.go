package repository

import (
	infraRedis "component-master/infra/redis"
	"context"
	"time"
)

type RedisRepository interface {
	GetValueByKey(ctx context.Context, key string) interface{}
	SetValueByKey(ctx context.Context, key string, value interface{}, ttl time.Duration)
}

type redisRepository struct {
	redis *infraRedis.RedisClient
}

func NewRedisRepository(redis *infraRedis.RedisClient) RedisRepository {
	return &redisRepository{redis: redis}
}

func (r *redisRepository) GetValueByKey(ctx context.Context, key string) interface{} {
	value, err := r.redis.Get(ctx, key)
	if err != nil {
		return nil
	}
	return value
}

func (r *redisRepository) SetValueByKey(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	err := r.redis.Set(ctx, key, value, ttl)
	if err != nil {
		return
	}
}
