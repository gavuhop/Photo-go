package service

import (
	"context"
	"fmt"
	"photo-go/internal/core"
	"photo-go/internal/database"
	"photo-go/pkg/logger"
	"photo-go/pkg/utils"
	"time"
)

type MediaService struct {
	VideoCore core.VideoProcessor
	ImageCore core.ImageProcessor
	Repo      database.MediaRepository
	Minio     *utils.MinioClient
}

func NewMediaService(v core.VideoProcessor, i core.ImageProcessor, r database.MediaRepository, m *utils.MinioClient) *MediaService {
	return &MediaService{
		VideoCore: v,
		ImageCore: i,
		Repo:      r,
		Minio:     m,
	}
}

// UploadAndProcessVideo upload, transcode, lưu DB
func (s *MediaService) UploadAndProcessVideo(ctx context.Context, filePath string, qualities []string) (string, error) {
	logger.Info("Start processing video: %s", filePath)
	// 1. Transcode HLS multi-quality ra thư mục tạm
	outputDir := "/tmp/hls_output" // TODO: random hóa tránh trùng
	if err := s.VideoCore.TranscodeToHLS(ctx, filePath, outputDir, qualities); err != nil {
		logger.Error(err, "TranscodeToHLS failed: %s", filePath)
		return "", err
	}
	logger.Info("TranscodeToHLS success: %s", filePath)
	// 2. Upload tất cả file .m3u8, .ts lên MinIO
	for _, q := range qualities {
		m3u8Path := fmt.Sprintf("%s/%s.m3u8", outputDir, q)
		objectName := fmt.Sprintf("hls/%d/%s.m3u8", time.Now().UnixNano(), q)
		logger.Info("Uploading m3u8 to Minio: %s", objectName)
		if err := s.Minio.Upload(ctx, objectName, m3u8Path); err != nil {
			logger.Error(err, "Minio upload failed: %s", m3u8Path)
			return "", err
		}
		logger.Info("Uploaded m3u8 to Minio: %s", objectName)
		// TODO: upload các file .ts tương ứng
	}
	// 3. Lưu DB
	media := &database.Media{
		Type:      "video",
		Path:      "hls/.../master.m3u8", // TODO: sinh master playlist
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	logger.Info("Saving media to DB: %s", media.Path)
	if err := s.Repo.Create(media); err != nil {
		logger.Error(err, "DB create media failed")
		return "", err
	}
	logger.Info("Media saved to DB: %s", media.Path)
	return media.Path, nil
}

// UploadAndProcessImage upload, xử lý, lưu DB
func (s *MediaService) UploadAndProcessImage(ctx context.Context, filePath string) error {
	// TODO: upload MinIO, gọi ImageCore.ProcessImage, lưu DB
	return nil
}
