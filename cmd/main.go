package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"photo-go/config"
	v1 "photo-go/internal/api/v1"
	"photo-go/internal/core"
	"photo-go/internal/database"
	"photo-go/internal/service"
	"photo-go/pkg/utils"
)

func main() {
	// Load config
	cfg := config.Settings

	// Init GORM
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := database.AutoMigrate(db); err != nil {
		panic(err)
	}

	// Init MinIO
	minioClient, err := utils.NewMinioClient(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		false, // useSSL
	)
	if err != nil {
		panic(err)
	}

	// Init core
	videoCore := core.NewFFMPEGVideoProcessor()
	imageCore := core.NewDefaultImageProcessor()

	// Init repository & service
	repo := database.NewGormMediaRepository(db)
	mediaService := service.NewMediaService(videoCore, imageCore, repo, minioClient)

	// Init Fiber
	app := fiber.New()

	// Register API
	handler := v1.NewMediaHandler(mediaService)
	handler.RegisterRoutes(app)

	// Start server
	app.Listen(cfg.Port)
}
