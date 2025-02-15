package service

import (
	"component-master/config"
	"component-master/infra/grpc/client"
	proto "component-master/proto/account"
	"component-master/util"
	"context"
	"log/slog"
	"sent-service/filewriter"
	"sent-service/worker"
	"time"
)

var ttl = time.Duration(24) * time.Hour

type AccountService interface {
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error)
	Start()
	Stop() error
}

type accountService struct {
	grpcclient client.AccountClient
	fw         *filewriter.FileWriter
	worker     *worker.BalanceWorker
}

func NewAccountService(client client.AccountClient, fw *filewriter.FileWriter, conf config.WorkerConfig) AccountService {
	return &accountService{
		grpcclient: client,
		fw:         fw,
		worker:     worker.NewBalanceWorker(conf, fw, client),
	}
}

func (s *accountService) Start() {
	s.worker.Start(util.ContextwithTimeout())
}

func (s *accountService) Stop() error {
	return s.worker.Stop()
}

func (s *accountService) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	slog.Info(util.StringPattern("Account service >> %d", req.AccountId))
	go s.grpcclient.CreateAccount(context.Background(), req)
	return &proto.CreateAccountResponse{Code: 0, Message: "Success"}, nil
}

func (s *accountService) BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	s.worker.Submit(util.ContextwithTimeout(), req)
	return &proto.BalanceChangeResponse{Code: 0, Message: "Success"}, nil
}
