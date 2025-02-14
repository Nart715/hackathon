package main

import (
	"component-master/config"
	grpcServer "component-master/infra/grpc/server"
	"component-master/util"
	"fmt"
	"log"
	"log/slog"
	"sent-service/filewriter"
	"sent-service/route"
	"sent-service/service"
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

	fw, err := filewriter.NewFileWriter("./log/transaction.log",
		filewriter.WithBufferSize(8192),           // 8KB buffer
		filewriter.WithFlushInterval(time.Second), // Flush every second
	)
	if err != nil {
		log.Fatalf("Failed to create file writer: %v", err)
	}
	defer fw.Close()

	go StartGrpcServer(conf, fw)
	StartHttpServer(conf, fw)
}

func StartHttpServer(cfg *config.Config, fw *filewriter.FileWriter) {
	route.InitHttpServer(cfg, fw)
}

func StartGrpcServer(cfg *config.Config, fw *filewriter.FileWriter) {
	transactionService := service.NewTransactionService(fw)
	transactionGrpcServer := grpcServer.InitTransactionGrpcServer(cfg.Server.Grpc, transactionService)
	slog.Info("start send grpc server")
	transactionGrpcServer.Start()
}
