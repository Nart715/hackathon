package service

import (
	"component-master/infra/grpc/client"
	proto "component-master/proto/account"
	"component-master/util"
	"context"
	"fmt"
	"log/slog"
	"sent-service/filewriter"
	"time"
)

var ttl = time.Duration(24) * time.Hour

type AccountService interface {
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error)
}

type accountService struct {
	grpcclient client.AccountClient
	fw         *filewriter.FileWriter
}

func NewAccountService(client client.AccountClient, fw *filewriter.FileWriter) AccountService {
	return &accountService{grpcclient: client, fw: fw}
}

func (s *accountService) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	slog.Info(util.StringPattern("Account service >> %d", req.AccountId))
	go s.grpcclient.CreateAccount(context.Background(), req)
	return &proto.CreateAccountResponse{Code: 0, Message: "Success"}, nil
}

func (s *accountService) BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	// transactionId, account id, amount, action, status
	line := util.StringPattern("%d_%d_%d_%d_%s", req.Tx, req.Ac, req.Am, req.Act, "PROCESSING")
	err := s.fw.WriteLine(line)
	if err != nil {
		slog.Error(fmt.Sprintf("write processing file has an error %s", err.Error()))
	}
	go s.grpcclient.BalanceChange(util.ContextwithTimeout(), req)
	return &proto.BalanceChangeResponse{Code: 0, Message: "Success"}, nil
}
