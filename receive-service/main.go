package main

import (
	mconfig "component-master/config"
	grpcServer "component-master/infra/grpc/server"
	infraRedis "component-master/infra/redis"
	"component-master/repository"
	"component-master/util"
	"flag"
	"fmt"
	"log/slog"
	"receive-service/biz"
	"receive-service/service"

	protoaccount "component-master/proto/account"
)

var (
	flagConf   string
	configPath string
	envFile    string
)

func init() {
	util.LoadEnv()
	flag.StringVar(&configPath, "config", "./config", "config file path")
	flag.StringVar(&envFile, "env", "", "env file path")
	flag.Parse()
}

func LoadConfig() *mconfig.Config {
	conf := &mconfig.Config{}
	err := mconfig.LoadConfig(configPath, envFile, conf)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to load config, %v", err))
	}
	return conf
}

func main() {
	conf := LoadConfig()
	slog.Info(fmt.Sprintf("config: %v", conf))
	StartBettingGrpcServer(conf)
}

func StartBettingGrpcServer(conf *mconfig.Config) {
	slog.Info("StartBettingGrpcServer start")
	accountService := initServices(conf)
	bettingGrpcServer := grpcServer.InitAccountGrpcServer(conf.Server.Grpc, accountService)
	bettingGrpcServer.Start()
}

func initServices(conf *mconfig.Config) protoaccount.AccountServiceServer {
	redis, err := infraRedis.NewInitRedisClient(&conf.Redis)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to init redis, %v", err))
		return nil
	}

	slog.Info(fmt.Sprintf("redis: %v", redis != nil))
	redisRepository := repository.NewRedisRepository(redis)
	accountBiz := biz.NewAccountBiz(redisRepository)
	accountService := service.NewAccountService(accountBiz)
	return accountService
}
