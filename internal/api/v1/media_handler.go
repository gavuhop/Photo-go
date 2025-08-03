package v1

import (
	"photo-go/pkg/logger"

	"github.com/gofiber/fiber/v3"
)

type MediaHandler struct {
	Service *MediaService
}

func NewMediaHandler(s *MediaService) *MediaHandler {
	return &MediaHandler{Service: s}
}

func (h *MediaHandler) RegisterRoutes(r fiber.Router) {
	r.Post("/media/upload", h.Upload)
	r.Get("/media/:id", h.Get)
	r.Get("/media/stream/:id", h.StreamHLS)
}

func (h *MediaHandler) Upload(c fiber.Ctx) error {
	logger.Info("Received upload request")
	file, err := c.FormFile("file")
	if err != nil {
		logger.Warn("Missing file in upload request")
		return c.Status(400).SendString("Missing file")
	}
	filePath := "/tmp/" + file.Filename
	if err := c.SaveFile(file, filePath); err != nil {
		logger.Error(err, "Save file error: %s", filePath)
		return c.Status(500).SendString("Save file error")
	}
	qualities := []string{"360p", "480p", "720p"}
	logger.Info("Processing video: %s", filePath)
	path, err := h.Service.UploadAndProcessVideo(c, filePath, qualities)
	if err != nil {
		logger.Error(err, "Error processing video: %s", filePath)
		return c.Status(500).SendString(err.Error())
	}
	logger.Info("Upload and process video success: %s", path)
	return c.JSON(fiber.Map{"stream_url": "/media/stream/" + path})
}

func (h *MediaHandler) Get(c fiber.Ctx) error {
	// TODO: Lấy media theo id, trả về link hoặc stream
	return c.SendStatus(501)
}

func (h *MediaHandler) StreamHLS(c fiber.Ctx) error {
	// id := c.Params("id") // Chưa dùng, comment lại để tránh lỗi
	// TODO: Lấy path m3u8 từ DB, redirect hoặc serve file MinIO
	return c.SendStatus(501)
}
