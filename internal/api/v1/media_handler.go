package v1

import (
	"fmt"
	"photo-go/internal/service"

	"github.com/gofiber/fiber/v3"
)

type MediaHandler struct {
	Service *service.MediaService
}

func NewMediaHandler(s *service.MediaService) *MediaHandler {
	return &MediaHandler{Service: s}
}

func (h *MediaHandler) RegisterRoutes(r fiber.Router) {
	r.Post("/media/upload", h.Upload)
	r.Get("/media/:id", h.Get)
	r.Get("/media/stream/:id", h.StreamHLS)
}

func (h *MediaHandler) Upload(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString("Missing file")
	}
	filePath := "/tmp/" + file.Filename
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).SendString("Save file error")
	}
	qualities := []string{"360p", "480p", "720p"}
	path, err := h.Service.UploadAndProcessVideo(c, filePath, qualities)
	fmt.Println(path, err)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
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
