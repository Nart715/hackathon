package route

import (
	"component-master/config"
	grpcClient "component-master/infra/grpc/client"
	"component-master/infra/http"
	"sent-service/handler"
	"sent-service/service"

	"github.com/gofiber/fiber/v2"
)

func InitHttpServer(conf *config.Config) {
	httpClient := http.HttpServer{
		AppName: "sent service",
		Conf:    &conf.Server.Http,
	}

	httpClient.InitHttpServer()

	v1 := httpClient.App().Group("/api/v1")
	SetupBettingRoute(v1, conf)
	httpClient.Start()
}

func SetupBettingRoute(r fiber.Router, conf *config.Config) {
	bettingClient := grpcClient.NewBettingClient(conf.GrpcClient)
	bettingService := service.NewBettingService(bettingClient)
	bettingHandler := handler.NewBettingHandler(bettingService)

	groupBetting := r.Group("/betting")
	POST(groupBetting, "", bettingHandler.CreateBetting)
}
