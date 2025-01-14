package server

import (
	"component-master/config"
	"component-master/proto/account"
)

func InitAccountGrpcServer(conf config.ServerInfo, accountServer account.AccountServiceServer) *GrpcServer {
	grpcServer := &GrpcServer{
		conf: &conf,
	}
	grpcServer.InitGrpcServer()
	account.RegisterAccountServiceServer(grpcServer.server, accountServer)
	return grpcServer
}
