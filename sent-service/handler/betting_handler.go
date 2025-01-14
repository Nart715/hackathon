package handler

import (
	"sent-service/service"

	"component-master/util"

	proto "component-master/proto/account"

	"github.com/gofiber/fiber/v2"
)

type AccountHandler interface {
	CreateAccount(c *fiber.Ctx) error
	BalanceChange(c *fiber.Ctx) error
}

type accountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler(accountService service.AccountService) AccountHandler {
	return &accountHandler{accountService: accountService}
}

func (b *accountHandler) CreateAccount(c *fiber.Ctx) error {
	ctx := util.ContextwithTimeout()
	req := &proto.CreateAccountRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	res, err := b.accountService.CreateAccount(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if res.GetCode() != 0 {
		return c.Status(fiber.StatusInternalServerError).SendString(res.GetMessage())
	}
	return c.Status(fiber.StatusOK).SendString("Account successfully created")
}

func (b *accountHandler) BalanceChange(c *fiber.Ctx) error {
	ctx := util.ContextwithTimeout()
	req := &proto.BalanceChangeRequest{}
	res, err := b.accountService.BalanceChange(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if res.GetCode() != 0 {
		return c.Status(fiber.StatusInternalServerError).SendString(res.GetMessage())
	}
	return c.Status(fiber.StatusOK).SendString("Balance successfully changed")
}
