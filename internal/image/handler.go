package image

import (
	"strings"

	"github.com/ahmadnaufal/openidea-paimonbank/internal/config"
	"github.com/ahmadnaufal/openidea-paimonbank/internal/model"
	"github.com/ahmadnaufal/openidea-paimonbank/pkg/jwt"
	"github.com/ahmadnaufal/openidea-paimonbank/pkg/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type imageHandler struct {
	s3Provider *s3.S3Provider
}

func NewImageHandler(s3prov *s3.S3Provider) imageHandler {
	return imageHandler{
		s3Provider: s3prov,
	}
}

func (h *imageHandler) RegisterRoute(r *fiber.App, jwtProvider jwt.JWTProvider) {
	imageGroup := r.Group("/v1/image")
	authMiddleware := jwtProvider.Middleware()

	imageGroup.Post("/", authMiddleware, h.UploadImage)
}

func (h *imageHandler) UploadImage(c *fiber.Ctx) error {
	// check for credentials
	_, err := jwt.GetLoggedInUser(c)
	if err != nil {
		return config.ErrRequestForbidden
	}

	fileReader, err := c.FormFile("file")
	if err != nil {
		return errors.Wrap(config.ErrInvalidUploadedFile, err.Error())
	}

	// check file size & extension
	fileSize := fileReader.Size
	if fileSize < 10*1024 || fileSize > 2*1024*1024 {
		return config.ErrInvalidFileSize
	}

	// check extension
	fp := strings.Split(fileReader.Filename, ".")
	if len(fp) < 2 || (fp[len(fp)-1] != "jpg" && fp[len(fp)-1] != "jpeg") {
		return config.ErrInvalidFileExtension
	}

	imgUrl, err := h.s3Provider.UploadImage(c.Context(), fileReader)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse{
		Message: "File uploaded successfully",
		Data: ImageUploadResponse{
			ImageURL: imgUrl,
		},
	})
}
