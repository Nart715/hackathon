package redis

import (
	"component-master/config"
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.ClusterClient
}

func NewInitRedisClient(rd *config.RedisConfig) (*RedisClient, error) {

	for _, addr := range rd.Clusters {
		slog.Info(fmt.Sprintf("Redis Clusters: %v", addr))
	}
	redisClusterOps := redis.ClusterOptions{
		Addrs:         rd.Clusters,
		ReadOnly:      rd.ReadOnly,
		RouteRandomly: rd.RouteRandomly,
		MaxRedirects:  rd.MaxRedirects,
		DialTimeout:   rd.DialTimeout,
		ReadTimeout:   rd.ReadTimeout,
		WriteTimeout:  rd.WriteTimeout,

		MinRetryBackoff: time.Millisecond * 100,
		MaxRetryBackoff: time.Second * 2,
		MinIdleConns:    rd.MinIdleConns,
		PoolSize:        rd.PoolSize,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Printf("Connected to Redis node: %s", cn.String())
			return nil
		},
	}

	client := redis.NewClusterClient(&redisClusterOps)
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("Ping to redis failed: %v", err))
		return nil, err
	}
	slog.Info(fmt.Sprintf("Ping to redis success: %s", pong))
	return &RedisClient{
		client: client,
	}, nil
}

func (r *RedisClient) GetRd() *redis.ClusterClient {
	return r.client
}

func (rc *RedisClient) Close() {
	if rc.client == nil {
		return
	}
	rc.client.Close()
}
