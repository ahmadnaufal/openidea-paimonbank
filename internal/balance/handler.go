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
	claims, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return config.ErrRequestForbidden
	}
	payload.UserID = claims.UserID

	if err := c.BodyParser(&payload); err != nil {
		return errors.Wrap(config.ErrMalformedRequest, err.Error())
	}

	if err := validation.Validate(payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	ctx := c.Context()

	// do addition
	transactionID := uuid.NewString()
	balanceEntity := BalanceHistory{
		ID:                      transactionID,
		UserID:                  payload.UserID,
		Currency:                payload.Currency,
		Balance:                 int(payload.AddedBalance),
		SourceBankAccountNumber: payload.SenderBankAccountNumber,
		SourceBankName:          payload.SenderBankName,
		TransferProofImg:        payload.TransferProofImg,
	}
	err = h.balanceRepo.AddBalance(ctx, nil, balanceEntity)
	if err != nil {
		return errors.Wrap(err, "AddBalance error")
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Message: "success",
		Data: BalanceHistoryResponse{
			TransactionID:    transactionID,
			Balance:          balanceEntity.Balance,
			Currency:         balanceEntity.Currency,
			TransferProofImg: balanceEntity.TransferProofImg,
			CreatedAt:        uint64(balanceEntity.CreatedAt.UnixMilli()),
			Source: BalanceSourceResponse{
				BankAccountNumber: balanceEntity.SourceBankAccountNumber,
				BankName:          balanceEntity.SourceBankName,
			},
		},
	})
}

func (h *balanceHandler) GetBalances(c *fiber.Ctx) error {
	claims, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return config.ErrRequestForbidden
	}

	ctx := c.Context()
	currencyBalances, err := h.balanceRepo.GetBalancePerCurrencies(ctx, claims.UserID)
	if err != nil {
		return errors.Wrap(err, "GetBalancePerCurrencies error")
	}

	responses := []CurrencyBalanceResponse{}
	for _, balance := range currencyBalances {
		responses = append(responses, CurrencyBalanceResponse(balance))
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Message: "success",
		Data:    responses,
	})
}

func (h *balanceHandler) GetBalanceHistory(c *fiber.Ctx) error {
	var payload GetBalanceHistoryRequest
	claims, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return config.ErrRequestForbidden
	}
	payload.UserID = claims.UserID

	if err := c.QueryParser(&payload); err != nil {
		return errors.Wrap(config.ErrMalformedRequest, err.Error())
	}

	if err := validation.Validate(payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	ctx := c.Context()
	balanceHistories, count, err := h.balanceRepo.GetBalanceHistory(ctx, payload)
	if err != nil {
		return errors.Wrap(err, "GetBalanceHistory error")
	}

	responses := []BalanceHistoryResponse{}
	for _, balanceEntity := range balanceHistories {
		responses = append(responses, BalanceHistoryResponse{
			TransactionID:    balanceEntity.ID,
			Balance:          balanceEntity.Balance,
			Currency:         balanceEntity.Currency,
			TransferProofImg: balanceEntity.TransferProofImg,
			CreatedAt:        uint64(balanceEntity.CreatedAt.UnixMilli()),
			Source: BalanceSourceResponse{
				BankAccountNumber: balanceEntity.SourceBankAccountNumber,
				BankName:          balanceEntity.SourceBankName,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Message: "success",
		Data:    responses,
		Meta: &model.ResponseMeta{
			Limit:  payload.Limit,
			Offset: payload.Offset,
			Total:  count,
		},
	})
}
