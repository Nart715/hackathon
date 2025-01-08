package main

import (
	mconfig "component-master/config"
	"component-master/infra/database"
	grpcServer "component-master/infra/grpc/server"
	infraRedis "component-master/infra/redis"
	"component-master/util"
	"flag"
	"fmt"
	"log/slog"
	"receive-service/biz"
	"receive-service/repository"
	"receive-service/service"

	protoaccount "component-master/proto/account"
	protobetting "component-master/proto/betting"
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
	bettingService, accountService := initServices(conf)
	bettingGrpcServer := grpcServer.InitBettingGrpcServer(conf.Server.Grpc, bettingService, accountService)
	bettingGrpcServer.Start()
}

func initServices(conf *mconfig.Config) (protobetting.BettingServiceServer, protoaccount.AccountServiceServer) {
	db, err := database.NewDataSource(&conf.Database)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to init database, %v", err))
		return nil, nil
	}
	redis, err := infraRedis.NewInitRedisClient(&conf.Redis)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to init redis, %v", err))
		return nil, nil
	}

	bettingRepository := repository.NewBettingRepository(db)
	redisRepository := repository.NewRedisRepository(redis)
	bettingBiz := biz.NewBettingBiz(bettingRepository, redisRepository)
	bettingService := service.NewBettingService(bettingBiz)

	accountRepository := repository.NewAccountRepository(db)
	accountBiz := biz.NewAccountBiz(accountRepository)
	accountService := service.NewAccountService(accountBiz)
	return bettingService, accountService
}
