package main

import (
	mconfig "component-master/config"
	mhandler "component-master/handler"
	grpcClient "component-master/infra/grpc/client"
	grpcServer "component-master/infra/grpc/server"
	"component-master/infra/kafka"
	infraRedis "component-master/infra/redis"
	"component-master/repository"
	"component-master/util"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"receive-service/biz"
	"receive-service/service"
	"receive-service/worker"
	"syscall"

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
	kafkaClient, err := kafka.NewKafkaClient(conf)
	if err != nil {
		slog.Error("Failed to create Kafka client", "error", err)
		return
	}

	accountService, accountBiz := initServices(conf, kafkaClient)
	bettingGrpcServer := grpcServer.InitAccountGrpcServer(conf.Server.Grpc, accountService)

	go SetupKafka(conf, accountBiz, kafkaClient)
	bettingGrpcServer.Start()
}

func initServices(conf *mconfig.Config, kafkaClient *kafka.KafkaClientConfig) (protoaccount.AccountServiceServer, biz.AccountBiz) {
	redis, err := infraRedis.NewInitRedisClient(&conf.Redis)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to init redis, %v", err))
		return nil, nil
	}

	slog.Info(fmt.Sprintf("redis: %v", redis != nil))
	transactionGrpcClient := grpcClient.NewTransactionClient(conf.GrpcClient)
	redisRepository := repository.NewRedisRepository(redis, transactionGrpcClient, conf.Redis.Channel)
	accountBiz := biz.NewAccountBiz(redisRepository)
	accountService := service.NewAccountService(accountBiz, kafkaClient)
	return accountService, accountBiz
}

func SetupKafka(cfg *mconfig.Config, accountBiz biz.AccountBiz, kafkaClient *kafka.KafkaClientConfig) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaHandler := mhandler.NewKafkaHandler(kafkaClient)
	workerProcessor := worker.NewBalanceListen(accountBiz)
	kafkaHandler.RegisterProcessor(cfg.KafkaTopic.BalanceChange, workerProcessor.Process)

	topics := []string{cfg.KafkaTopic.BalanceChange}

	if err := kafkaClient.Consume(ctx, topics, kafkaHandler.Handle); err != nil {
		slog.Error("Failed to start consumer", "error", err)
		return err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()
	if err := kafkaClient.Close(); err != nil {
		slog.Error("Error closing Kafka client", "error", err)
		return err
	}

	return nil
}
