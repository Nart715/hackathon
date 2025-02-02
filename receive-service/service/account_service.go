package service

import (
	"component-master/infra/kafka"
	proto "component-master/proto/account"
	"context"
	"encoding/json"
	"receive-service/biz"
	"time"
)

type accountService struct {
	proto.UnimplementedAccountServiceServer
	accountBiz  biz.AccountBiz
	kafkaClient *kafka.KafkaClientConfig
}

func NewAccountService(accountBiz biz.AccountBiz, kafkaClient *kafka.KafkaClientConfig) proto.AccountServiceServer {
	return &accountService{accountBiz: accountBiz, kafkaClient: kafkaClient}
}

func (s accountService) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	_, err := s.accountBiz.CreateAccount(ctx, req)
	if err != nil {
		return &proto.CreateAccountResponse{
			Code:    -1,
			Message: err.Error(),
		}, nil
	}
	return &proto.CreateAccountResponse{
		Code:    0,
		Message: "success",
	}, nil
}

func (s accountService) BalanceChange(ctx context.Context, req *proto.BalanceChangeRequest) (*proto.BalanceChangeResponse, error) {
	topic := s.kafkaClient.GetConfig().KafkaTopic.BalanceChange
	data, err := json.Marshal(req)
	if err != nil {
		return &proto.BalanceChangeResponse{
			Code:    -1,
			Message: err.Error(),
		}, nil
	}
	s.kafkaClient.Produce(ctx, topic, keyGen(), data)
	return &proto.BalanceChangeResponse{
		Code:    0,
		Message: "success",
	}, nil
}

func keyGen() []byte {
	resp := []byte{byte(time.Now().UnixMicro())}
	return resp
}
