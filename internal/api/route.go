package api

import (
	v1 "photo-go/internal/api/v1"
	"photo-go/internal/core"
	"photo-go/pkg/utils"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Đăng ký tất cả route version 1 vào app
func RegisterV1Routes(app *fiber.App, db *gorm.DB, videoCore core.VideoProcessor, imageCore core.ImageProcessor, minioClient *utils.MinioClient) {
	repo := v1.NewGormMediaRepository(db)
	mediaService := v1.NewMediaService(videoCore, imageCore, repo, minioClient)
	handler := v1.NewMediaHandler(mediaService)
	v1Group := app.Group("/v1")
	handler.RegisterRoutes(v1Group)
}
