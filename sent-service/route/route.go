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

	v1 := httpClient.App().Group("/api/v1/player")
	SetupRoute(v1, conf)
	httpClient.Start()
}

func SetupRoute(r fiber.Router, conf *config.Config) {
	accountClient := grpcClient.NewAccountClient(conf.GrpcClient)
	accountService := service.NewAccountService(accountClient)
	accountHandler := handler.NewAccountHandler(accountService)

	groupCreatedAccount := r.Group("/created-account")
	POST(groupCreatedAccount, "", accountHandler.CreateAccount)

	// both of deposit and betting differs by the action
	groupBalanceChange := r.Group("/balance-change")
	POST(groupBalanceChange, "", accountHandler.BalanceChange)
}
