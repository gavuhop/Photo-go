package utils

import (
	"context"

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
	_, err := m.Client.FPutObject(ctx, m.Bucket, objectName, filePath, minio.PutObjectOptions{})
	return err
}

func (m *MinioClient) Download(ctx context.Context, objectName, filePath string) error {
	return m.Client.FGetObject(ctx, m.Bucket, objectName, filePath, minio.GetObjectOptions{})
}
