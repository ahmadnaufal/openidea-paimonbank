package balance

import (
	"github.com/ahmadnaufal/openidea-paimonbank/internal/config"
	"github.com/ahmadnaufal/openidea-paimonbank/internal/model"
	"github.com/ahmadnaufal/openidea-paimonbank/pkg/jwt"
	"github.com/ahmadnaufal/openidea-paimonbank/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type balanceHandler struct {
	balanceRepo *balanceRepo
}

type BalanceHandlerConfig struct {
	BalanceRepo *balanceRepo
}

func NewBalance(cfg BalanceHandlerConfig) balanceHandler {
	return balanceHandler{balanceRepo: cfg.BalanceRepo}
}

func (h *balanceHandler) RegisterRoute(r *fiber.App, jwtProvider jwt.JWTProvider) {
	imageGroup := r.Group("/v1/balance")
	authMiddleware := jwtProvider.Middleware()

	imageGroup.Post("/", authMiddleware, h.AddBalance)
	imageGroup.Get("/", authMiddleware, h.GetBalances)
	imageGroup.Get("/history", authMiddleware, h.GetBalanceHistory)
}

func (h *balanceHandler) AddBalance(c *fiber.Ctx) error {
	var payload AddBalanceRequest
	if err := c.BodyParser(&payload); err != nil {
		return errors.Wrap(config.ErrMalformedRequest, err.Error())
	}

	if err := validation.Validate(payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	transactionID := uuid.NewString()
	// do addition

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Message: "success",
		Data: BalanceHistoryResponse{
			TransactionID:    transactionID,
			Balance:          int(payload.AddedBalance),
			Currency:         payload.Currency,
			TransferProofImg: payload.TransferProofImg,
			Source: BalanceSourceResponse{
				BankAccountNumber: payload.SenderBankAccountNumber,
				BankName:          payload.SenderBankName,
			},
		},
	})
}

func (h *balanceHandler) GetBalances(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Message: "success",
		Data:    CurrencyBalanceResponse{},
	})
}

func (h *balanceHandler) GetBalanceHistory(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Message: "success",
		Data:    []BalanceHistoryResponse{},
		Meta:    &model.ResponseMeta{},
	})
}
