package service

import (
	"component-master/infra/grpc/client"
	proto "component-master/proto/account"
	"component-master/util"
	"context"
	"fmt"
	"log/slog"
	"sent-service/filewriter"
	"sync"
	"time"
)

// Số lần thử lại nếu request gRPC thất bại
const maxRetries = 3

type AccountService interface {
	CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error)
	BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error)
}

type accountService struct {
	grpcClientPool *sync.Pool
	fw             *filewriter.FileWriter
}

func NewAccountService(clientFactory func() client.AccountClient, fw *filewriter.FileWriter) AccountService {
	// Tạo pool chứa gRPC client
	clientPool := &sync.Pool{
		New: func() interface{} {
			return clientFactory()
		},
	}
	return &accountService{grpcClientPool: clientPool, fw: fw}
}

// Hàm lấy client từ pool
func (s *accountService) getGRPCClient() client.AccountClient {
	return s.grpcClientPool.Get().(client.AccountClient)
}

// Trả client về pool
func (s *accountService) releaseGRPCClient(client client.AccountClient) {
	s.grpcClientPool.Put(client)
}

func (s *accountService) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	slog.Info(util.StringPattern("Account service >> %d", req.AccountId))
	go func() {
		client := s.getGRPCClient()
		defer s.releaseGRPCClient(client)
		client.CreateAccount(context.Background(), req)
	}()
	return &proto.CreateAccountResponse{Code: 0, Message: "Success"}, nil
}

func (s *accountService) BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	// Ghi vào file buffer để tránh blocking IO
	line := fmt.Sprintf("%d_%d_%d_%d_%s", req.Tx, req.Ac, req.Am, req.Act, "PROCESSING")
	go func() {
		if err := s.fw.WriteLine(line); err != nil {
			slog.Error(fmt.Sprintf("write processing file error: %s", err.Error()))
		}
	}()

	// Xử lý gRPC trong worker pool
	go func() {
		client := s.getGRPCClient()
		defer s.releaseGRPCClient(client)

		var err error
		for i := 0; i < maxRetries; i++ {
			_, err = client.BalanceChange(util.ContextwithTimeout(), req)
			if err == nil {
				return
			}
			time.Sleep(50 * time.Millisecond) // Thử lại sau 50ms nếu lỗi
		}
		slog.Error(fmt.Sprintf("BalanceChange failed after retries: %v", err))
	}()

	return &proto.BalanceChangeResponse{Code: 0, Message: "Success"}, nil
}
