package service

import (
	proto "component-master/proto/transaction"
	"component-master/util"
	"context"
	"log/slog"
	"sent-service/filewriter"
)

type transactionService struct {
	proto.UnimplementedTransactionServiceServer
	fw *filewriter.FileWriter
}

func NewTransactionService(fw *filewriter.FileWriter) proto.TransactionServiceServer {
	return &transactionService{fw: fw}
}

func (r *transactionService) UpdateTransaction(ctx context.Context, req *proto.TransactionUpdate) (*proto.Empty, error) {
	slog.Info(util.StringPattern("update transaction %d status from receive", req.Tx))
	if req == nil {
		slog.Error("Update transaction is failure by nil")
		return &proto.Empty{}, nil
	}

	if req.Tx == 0 {
		slog.Error("Update transaction is failure by zero")
		return &proto.Empty{}, nil
	}
	if req.Status == "" {
		slog.Error("Update transaction is failure by empty")
		return &proto.Empty{}, nil
	}
	if r.fw == nil {
		slog.Error("Update transaction is failure by write file is nil")
		return &proto.Empty{}, nil
	}

	// transaction.internal transaction, account, balance before, change, balance after, type, status
	line := util.StringPattern("%d.%d_%d_%d_%d_%d_%d_%s", req.Tx, req.Itx, req.Ac, req.Bab, req.Change, req.Baf, req.Type, req.Status)
	err := r.fw.WriteLine(line)
	if err != nil {
		slog.Error("Update transaction status is failure")
		return &proto.Empty{}, err
	}
	return &proto.Empty{}, nil
}
