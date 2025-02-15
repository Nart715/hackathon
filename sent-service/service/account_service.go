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

var ttl = 24 * time.Hour

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

	// Sử dụng context với timeout hợp lý
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := s.grpcclient.CreateAccount(ctx, req)
	if err != nil {
		slog.Error(fmt.Sprintf("gRPC CreateAccount failed: %s", err.Error()))
		return &proto.CreateAccountResponse{Code: -1, Message: "Failed"}, err
	}

	return resp, nil
}

func (s *accountService) BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	// Ghi log trạng thái xử lý
	line := util.StringPattern("%d_%d_%d_%d_%s", req.Tx, req.Ac, req.Am, req.Act, "PROCESSING")
	err := s.fw.WriteLine(line)
	if err != nil {
		slog.Error(fmt.Sprintf("write processing file has an error %s", err.Error()))
	}

	// Sử dụng context với timeout hợp lý
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := s.grpcclient.BalanceChange(ctx, req)
	if err != nil {
		slog.Error(fmt.Sprintf("gRPC BalanceChange failed: %s", err.Error()))
		return &proto.BalanceChangeResponse{Code: -1, Message: "Failed"}, err
	}

	return resp, nil
}
