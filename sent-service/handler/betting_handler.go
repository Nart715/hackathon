package handler

import (
	"sent-service/service"
	"time"

	"component-master/util"

	proto "component-master/proto/message"

	"github.com/gofiber/fiber/v2"
)

type BettingHandler interface {
	CreateBetting(c *fiber.Ctx) error
}

type bettingHandler struct {
	bettingService service.BettingService
}

func NewBettingHandler(bettingService service.BettingService) BettingHandler {
	return &bettingHandler{bettingService: bettingService}
}

func (b *bettingHandler) CreateBetting(c *fiber.Ctx) error {
	ctx := util.ContextwithTimeout()
	// TODO: scan request from api to here
	b.bettingService.CreateBetting(ctx, &proto.BettingMessageRequest{
		AccountId:     "1",
		Amount:        100,
		BetId:         100,
		TransactionId: 100,
		Timestamp:     time.Now().UnixMilli(),
	})
	return c.Status(fiber.StatusOK).SendString("Betting successfully created")
}
