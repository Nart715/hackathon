package main

import (
	"component-master/config"
	"component-master/util"
	"fmt"
	"log/slog"
	"sent-service/route"
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

	route.InitHttpServer(conf)
}
