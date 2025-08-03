package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"photo-go/config"
	v1 "photo-go/internal/api/v1"
	"photo-go/internal/core"
	"photo-go/internal/database"
	"photo-go/internal/service"
	"photo-go/pkg/logger"
	"photo-go/pkg/utils"
)

func main() {
	logger.Info("Starting application")
	// Load config
	cfg := config.Settings

	// Init GORM
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal(err, "Failed to connect to DB")
	}
	logger.Info("Connected to DB")
	if err := database.AutoMigrate(db); err != nil {
		logger.Fatal(err, "DB migration failed")
	}
	logger.Info("DB migration success")

	// Init MinIO
	minioClient, err := utils.NewMinioClient(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		false, // useSSL
	)
	if err != nil {
		logger.Fatal(err, "Failed to init Minio client")
	}
	logger.Info("Minio client initialized")

	// Init core
	videoCore := core.NewFFMPEGVideoProcessor()
	imageCore := core.NewDefaultImageProcessor()
	logger.Info("Core processors initialized")

	// Init repository & service
	repo := database.NewGormMediaRepository(db)
	mediaService := service.NewMediaService(videoCore, imageCore, repo, minioClient)
	logger.Info("Media service initialized")

	// Init Fiber
	app := fiber.New()
	logger.Info("Fiber app initialized")

	// CORS middleware với cấu hình từ config (truyền slice trực tiếp)
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "Cache-Control"},
	}))

	// Register API
	handler := v1.NewMediaHandler(mediaService)
	handler.RegisterRoutes(app)
	logger.Info("API routes registered")

	// Start server
	logger.Info("Starting server on port %s", cfg.Port)
	app.Listen(cfg.Port)
}
