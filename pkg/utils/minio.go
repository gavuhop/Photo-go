package utils

import (
	"context"

	"photo-go/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client *minio.Client
	Bucket string
}

func NewMinioClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinioClient{Client: cli, Bucket: bucket}, nil
}

func (m *MinioClient) Upload(ctx context.Context, objectName, filePath string) error {
	logger.Info("Uploading to Minio: %s from %s", objectName, filePath)
	_, err := m.Client.FPutObject(ctx, m.Bucket, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		logger.Error(err, "Minio upload failed: %s", objectName)
	} else {
		logger.Info("Minio upload success: %s", objectName)
	}
	return err
}

func (m *MinioClient) Download(ctx context.Context, objectName, filePath string) error {
	logger.Info("Downloading from Minio: %s to %s", objectName, filePath)
	err := m.Client.FGetObject(ctx, m.Bucket, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		logger.Error(err, "Minio download failed: %s", objectName)
	} else {
		logger.Info("Minio download success: %s", objectName)
	}
	return err
}
