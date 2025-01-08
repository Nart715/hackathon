package server

import (
	"component-master/config"
	"component-master/proto/account"
	"component-master/proto/betting"
)

func InitBettingGrpcServer(conf config.ServerInfo, bettingServer betting.BettingServiceServer, accountServer account.AccountServiceServer) *GrpcServer {
	grpcServer := &GrpcServer{
		conf: &conf,
	}
	grpcServer.InitGrpcServer()
	betting.RegisterBettingServiceServer(grpcServer.server, bettingServer)
	account.RegisterAccountServiceServer(grpcServer.server, accountServer)
	return grpcServer
}
