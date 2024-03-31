package config

import (
	"net/http"

	"github.com/ahmadnaufal/openidea-paimonbank/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

var (
	ErrMalformedRequest     = fiber.NewError(http.StatusBadRequest, "request malformed")
	ErrCredentialExists     = fiber.NewError(http.StatusConflict, "credential already used")
	ErrWrongPassword        = fiber.NewError(http.StatusBadRequest, "wrong password entered")
	ErrRequestForbidden     = fiber.NewError(http.StatusForbidden, "request forbidden")
	ErrInsufficientBalance  = fiber.NewError(http.StatusBadRequest, "insufficient balance in currency")
	ErrUserNotFound         = fiber.NewError(http.StatusNotFound, "user with the specified credential not found")
	ErrPostNotFound         = fiber.NewError(http.StatusNotFound, "post not found")
	ErrInvalidUploadedFile  = fiber.NewError(http.StatusBadRequest, "invalid uploaded file")
	ErrInvalidFileSize      = fiber.NewError(http.StatusBadRequest, "invalid file size")
	ErrInvalidFileExtension = fiber.NewError(http.StatusBadRequest, "invalid file extension")
)

func DefaultErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Status code defaults to 500
		code := fiber.StatusInternalServerError
		message := "internal server error"

		// Retrieve the custom status code & message if it's a *fiber.Error
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
			message = e.Message
		}

		// Return status code with error message
		return c.Status(code).JSON(model.ErrorResponse{
			Message: message,
		})
	}
}
