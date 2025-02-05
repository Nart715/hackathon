package workerkafka

import (
	proto "component-master/proto/account"
	"component-master/util"
	"context"
	"encoding/json"
	"log/slog"
	"receive-service/biz"
)

type BalanceListenProcessor interface {
	Process(context.Context, []byte) error
}

type balanceListenProcessor struct {
	accountBiz biz.AccountBiz
}

func NewBalanceListen(accountBiz biz.AccountBiz) BalanceListenProcessor {
	return &balanceListenProcessor{accountBiz: accountBiz}
}

func (r *balanceListenProcessor) Process(ctx context.Context, data []byte) error {
	slog.Info(string(data))
	req := &proto.BalanceChangeRequest{}
	err := json.Unmarshal(data, req)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	r.accountBiz.BalanceChange(util.ContextwithTimeout(), req)
	return nil
}
