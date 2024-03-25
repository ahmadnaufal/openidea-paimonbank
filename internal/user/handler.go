package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/ahmadnaufal/openidea-paimonbank/internal/config"
	"github.com/ahmadnaufal/openidea-paimonbank/internal/model"
	"github.com/ahmadnaufal/openidea-paimonbank/pkg/jwt"
	"github.com/ahmadnaufal/openidea-paimonbank/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type userHandler struct {
	userRepo    *UserRepo
	jwtProvider *jwt.JWTProvider
	saltCost    int
}

type UserHandlerConfig struct {
	UserRepo    *UserRepo
	JwtProvider *jwt.JWTProvider
	SaltCost    int
}

func NewUserHandler(cfg UserHandlerConfig) userHandler {
	return userHandler{
		userRepo:    cfg.UserRepo,
		jwtProvider: cfg.JwtProvider,
		saltCost:    cfg.SaltCost,
	}
}

func (h *userHandler) RegisterRoute(r *fiber.App, _ jwt.JWTProvider) {
	userGroup := r.Group("/v1/user")

	userGroup.Post("/register", h.RegisterUser)
	userGroup.Post("/login", h.Authenticate)
}

func (h *userHandler) RegisterUser(c *fiber.Ctx) error {
	var payload RegisterUserRequest
	if err := c.BodyParser(&payload); err != nil {
		return errors.Wrap(config.ErrMalformedRequest, err.Error())
	}

	if err := validation.Validate(payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// find existing user by credentials
	user, accessToken, err := h.createUser(c.Context(), payload)
	if err != nil {
		return errors.Wrap(err, "create user error")
	}

	return c.Status(fiber.StatusCreated).JSON(model.DataResponse{
		Message: "User registered successfully",
		Data: UserResponse{
			Email:       user.Email,
			Name:        user.Name,
			AccessToken: accessToken,
		},
	})
}

func (h *userHandler) createUser(ctx context.Context, payload RegisterUserRequest) (User, string, error) {
	_, err := h.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil && err != sql.ErrNoRows {
		return User{}, "", errors.Wrap(err, "GetUserByCredential error")
	}
	if err == nil {
		// user already exists
		return User{}, "", config.ErrCredentialExists
	}

	// hash the password first using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), h.saltCost)
	if err != nil {
		return User{}, "", err
	}

	user := User{
		ID:       uuid.NewString(),
		Name:     payload.Name,
		Email:    payload.Email,
		Password: string(hashedPassword),
	}
	err = h.userRepo.CreateUser(ctx, user)
	if err != nil {
		return user, "", err
	}

	// generate JWT
	accessToken, err := h.generateAccessTokenFromUser(user)
	if err != nil {
		return user, "", err
	}

	return user, accessToken, nil
}

func (h *userHandler) Authenticate(c *fiber.Ctx) error {
	var payload AuthenticateRequest
	if err := c.BodyParser(&payload); err != nil {
		return errors.Wrap(config.ErrMalformedRequest, err.Error())
	}

	if err := validation.Validate(payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, accessToken, err := h.authenticateUser(c.Context(), payload)
	if err != nil {
		return errors.Wrap(err, "create user error")
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Message: "User logged successfully",
		Data: UserResponse{
			Email:       user.Email,
			Name:        user.Name,
			AccessToken: accessToken,
		},
	})
}

func (h *userHandler) authenticateUser(ctx context.Context, payload AuthenticateRequest) (User, string, error) {
	user, err := h.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, "", config.ErrUserNotFound
		}

		return user, "", errors.Wrap(err, "GetUserByCredential error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return user, "", config.ErrWrongPassword
	}

	// generate JWT
	accessToken, err := h.generateAccessTokenFromUser(user)
	if err != nil {
		return user, "", errors.Wrap(err, "generateAccessToken error")
	}

	return user, accessToken, nil
}

func (h *userHandler) generateAccessTokenFromUser(user User) (string, error) {
	claims := jwt.BuildJWTClaims(jwt.JWTUser{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}, 8*time.Hour)

	accessToken, err := h.jwtProvider.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return accessToken, err
}
