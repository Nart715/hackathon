package main

import (
	"component-master/config"
	"component-master/infra/database"
	"component-master/infra/nat"
	"component-master/infra/redis"
	"component-master/util"
	"fmt"
	"log/slog"
)

func init() {
	util.LoadEnv()
}

func main() {
	conf := &config.Config{}
	err := config.LoadConfig("config", "", conf)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to load config, %v", err))
		return
	}

	dbConf := conf.Database
	fmt.Println(dbConf.DriverName)

	db, err := database.NewDataSource(&dbConf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rd, err := redis.NewInitRedisClient(&conf.Redis)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to init redis, %v", err))
		return
	}
	slog.Info(fmt.Sprintf("REDIS CONNECT: %v", rd != nil))

	defer db.Close()

	slog.Info(fmt.Sprintf("NATS CONNECT %s:%d", conf.Nats.Host, conf.Nats.Port))
	natsClient, err := nat.NewNatsConnectClient(&conf.Nats)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to init nats, %v", err))
		return
	}
	slog.Info(fmt.Sprintf("NATS CONNECT: %v", natsClient != nil))
}
