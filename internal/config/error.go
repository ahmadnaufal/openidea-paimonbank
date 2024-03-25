package config

import (
	"net/http"

	"github.com/ahmadnaufal/openidea-paimonbank/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

var (
	ErrMalformedRequest       = fiber.NewError(http.StatusBadRequest, "request malformed")
	ErrCredentialExists       = fiber.NewError(http.StatusConflict, "credential already used")
	ErrWrongPassword          = fiber.NewError(http.StatusBadRequest, "wrong password entered")
	ErrRequestForbidden       = fiber.NewError(http.StatusForbidden, "request forbidden")
	ErrCannotChangeCredential = fiber.NewError(http.StatusBadRequest, "cannot change email or phone from link email/phone API")
	ErrUserNotFound           = fiber.NewError(http.StatusNotFound, "user with the specified credential not found")
	ErrPostNotFound           = fiber.NewError(http.StatusNotFound, "post not found")
	ErrPostCreatorIsNotFriend = fiber.NewError(http.StatusBadRequest, "you cannot comment a post which author is not on your friend list")
	ErrTargetUserIDEmpty      = fiber.NewError(http.StatusBadRequest, "user ID is empty")
	ErrSelfAddFriend          = fiber.NewError(http.StatusBadRequest, "cannot add yourself as a new friend")
	ErrFriendAlreadyAdded     = fiber.NewError(http.StatusBadRequest, "user already added as friend")
	ErrUserIsNotAFriend       = fiber.NewError(http.StatusBadRequest, "user is not a friend")
	ErrInvalidUploadedFile    = fiber.NewError(http.StatusBadRequest, "invalid uploaded file")
	ErrInvalidFileSize        = fiber.NewError(http.StatusBadRequest, "invalid file size")
	ErrInvalidFileExtension   = fiber.NewError(http.StatusBadRequest, "invalid file extension")
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
