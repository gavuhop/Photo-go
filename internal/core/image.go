package core

import (
	"context"
)

// ImageProcessor định nghĩa interface xử lý ảnh

type ImageProcessor interface {
	ProcessImage(ctx context.Context, inputPath, outputPath string) error
}

type DefaultImageProcessor struct{}

func NewDefaultImageProcessor() *DefaultImageProcessor {
	return &DefaultImageProcessor{}
}

func (p *DefaultImageProcessor) ProcessImage(ctx context.Context, inputPath, outputPath string) error {
	// TODO: Xử lý ảnh (resize, crop, ...)
	return nil
}
