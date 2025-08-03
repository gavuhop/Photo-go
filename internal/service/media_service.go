package service

import (
	"context"
	"fmt"
	"photo-go/internal/core"
	"photo-go/internal/database"
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
	// 1. Transcode HLS multi-quality ra thư mục tạm
	outputDir := "/tmp/hls_output" // TODO: random hóa tránh trùng
	if err := s.VideoCore.TranscodeToHLS(ctx, filePath, outputDir, qualities); err != nil {
		return "", err
	}
	// 2. Upload tất cả file .m3u8, .ts lên MinIO
	for _, q := range qualities {
		m3u8Path := fmt.Sprintf("%s/%s.m3u8", outputDir, q)
		objectName := fmt.Sprintf("hls/%d/%s.m3u8", time.Now().UnixNano(), q)
		if err := s.Minio.Upload(ctx, objectName, m3u8Path); err != nil {
			return "", err
		}
		// TODO: upload các file .ts tương ứng
	}
	// 3. Lưu DB
	media := &database.Media{
		Type:      "video",
		Path:      "hls/.../master.m3u8", // TODO: sinh master playlist
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	if err := s.Repo.Create(media); err != nil {
		return "", err
	}
	return media.Path, nil
}

// UploadAndProcessImage upload, xử lý, lưu DB
func (s *MediaService) UploadAndProcessImage(ctx context.Context, filePath string) error {
	// TODO: upload MinIO, gọi ImageCore.ProcessImage, lưu DB
	return nil
}
