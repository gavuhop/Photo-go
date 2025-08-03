package v1

import (
	"photo-go/internal/core"
	"photo-go/pkg/utils"
)

type MediaService struct {
	VideoCore core.VideoProcessor
	ImageCore core.ImageProcessor
	Repo      MediaRepository
	Minio     *utils.MinioClient
}
