package route

import (
	"component-master/config"
	grpcClient "component-master/infra/grpc/client"
	"component-master/infra/http"
	"sent-service/filewriter"
	"sent-service/handler"
	"sent-service/service"

	"github.com/gofiber/fiber/v2"
)

func InitHttpServer(conf *config.Config, fw *filewriter.FileWriter) {
	httpClient := http.HttpServer{
		AppName: "sent service",
		Conf:    &conf.Server.Http,
	}

	httpClient.InitHttpServer()

	accountClient := grpcClient.NewAccountClient(conf.GrpcClient)
	accountService := service.NewAccountService(accountClient, fw, conf.Worker)
	accountHandler := handler.NewAccountHandler(accountService)

	v1 := httpClient.App().Group("/api/v1/player")
	SetupRoute(v1, conf, fw, accountHandler)
	httpClient.Start()
}

func SetupRoute(r fiber.Router, conf *config.Config, fw *filewriter.FileWriter, accountHandler handler.AccountHandler) {

	groupCreatedAccount := r.Group("/created-account")
	POST(groupCreatedAccount, "", accountHandler.CreateAccount)

	// both of deposit and betting differs by the action
	groupBalanceChange := r.Group("/balance-change")
	POST(groupBalanceChange, "", accountHandler.BalanceChange)
}
