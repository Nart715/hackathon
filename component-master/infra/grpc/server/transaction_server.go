package server

import (
	"component-master/config"
	proto "component-master/proto/transaction"
)

func InitTransactionGrpcServer(conf config.ServerInfo, transactionServer proto.TransactionServiceServer) *GrpcServer {
	grpcServer := &GrpcServer{
		conf: &conf,
	}
	grpcServer.InitGrpcServer()
	proto.RegisterTransactionServiceServer(grpcServer.server, transactionServer)
	return grpcServer
}
