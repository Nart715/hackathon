package main

import (
	"component-master/config"
	infraRedis "component-master/infra/redis"
	proto "component-master/proto/account"
	"component-master/repository"
	"component-master/util"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sent-service/route"
	"strconv"
	"time"
)

func init() {
	util.LoadEnv()
}

func main() {
	slog.Info("Hello World")
	conf := &config.Config{}
	err := config.LoadConfig("config", "", conf)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to load config, %v", err))
		return
	}

	redisRepo := InitRedisRepo(conf)
	defer redisRepo.Close()
	go InitRedisInfra(conf, redisRepo)
	route.InitHttpServer(conf, redisRepo)
}

func InitRedisRepo(cf *config.Config) repository.RedisRepository {
	redis, err := infraRedis.NewInitRedisClient(&cf.Redis)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to init redis, %v", err))
	}
	slog.Info(cf.Redis.Channel)
	redisRepository := repository.NewRedisRepository(redis, cf.Redis.Channel)
	return redisRepository
}

func InitRedisInfra(cf *config.Config, redisRepo repository.RedisRepository) {
	redisRepo.SubscribeMessage(context.Background(), process)

}

func process(message string, redisRepo repository.RedisRepository) {
	slog.Info("Message receive " + message)
	clazzz := &proto.AccountData{}
	if err := json.Unmarshal([]byte(message), clazzz); err != nil || clazzz.GetTx() <= 0 {
		return
	}
	key := strconv.Itoa(int(clazzz.GetTx()))
	data := redisRepo.GetRedis().Get(context.Background(), key).Val()
	if empty(data) {
		slog.Warn("Transaction not valid " + key)
		return
	}
	slog.Info("Update redis by key " + key)
	redisRepo.GetRedis().Set(context.Background(), key, "success", time.Duration(24)*time.Hour)
}

func empty(data string) bool {
	return len(data) <= 0
}
