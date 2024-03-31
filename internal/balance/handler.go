package balance

import (
	"context"
	"strings"

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
	authMiddleware := jwtProvider.Middleware()

	balanceGroup := r.Group("/v1/balance")
	balanceGroup.Post("/", authMiddleware, h.AddBalance)
	balanceGroup.Get("/", authMiddleware, h.GetBalances)
	balanceGroup.Get("/history", authMiddleware, h.GetBalanceHistory)

	transactionGroup := r.Group("/v1/transaction")
	transactionGroup.Post("/", authMiddleware, h.CreateTransaction)
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
	currencyBalances, err := h.balanceRepo.GetBalancePerCurrencies(ctx, claims.UserID, "")
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

func (h *balanceHandler) CreateTransaction(c *fiber.Ctx) error {
	var payload CreateTransactionRequest
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

	balanceEntity, err := h.createTransaction(c.Context(), payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Message: "success",
		Data: BalanceHistoryResponse{
			TransactionID:    balanceEntity.ID,
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

func (h *balanceHandler) createTransaction(ctx context.Context, payload CreateTransactionRequest) (BalanceHistory, error) {
	// first, validate if the budget does exist
	currencyBudgets, err := h.balanceRepo.GetBalancePerCurrencies(ctx, payload.UserID, payload.FromCurrency)
	if err != nil {
		return BalanceHistory{}, errors.Wrap(err, "GetBalancePerCurrencies error")
	}
	if len(currencyBudgets) != 1 || (currencyBudgets[0].Balance < int(payload.Balances)) {
		return BalanceHistory{}, config.ErrInsufficientBalance
	}

	normalizedCurrency := strings.ToUpper(payload.FromCurrency)
	deductedBalance := int(payload.Balances) * -1
	// then, save the balance history
	transactionID := uuid.NewString()
	balanceEntity := BalanceHistory{
		ID:                      transactionID,
		UserID:                  payload.UserID,
		Balance:                 deductedBalance,
		Currency:                normalizedCurrency,
		TransferProofImg:        "",
		SourceBankAccountNumber: payload.RecipientBankAccountNumber,
		SourceBankName:          payload.RecipientBankName,
	}
	err = h.balanceRepo.AddBalance(ctx, nil, balanceEntity)
	if err != nil {
		return BalanceHistory{}, errors.Wrap(err, "AddBalance error")
	}

	return balanceEntity, nil
}
