package minio_storage

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/config"
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client *minio.Client
	conf   *config.MinioConfig
}

// NewMinioClient создает новый MinioClient и инициализирует клиента MinIO
func NewMinioClient(conf *config.MinioConfig) (*MinioClient, error) {
	// Создаем нового клиента MinIO
	client, err := minio.New(conf.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания клиента MinIO: %w", err)
	}

	return &MinioClient{
		client: client,
		conf:   conf,
	}, nil
}

// ensureBucketExists проверяет, существует ли бакет, и создает его, если нет
func (m *MinioClient) ensureBucketExists(ctx context.Context, bucketName string, region string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("ошибка при проверке существования бакета: %w", err)
	}
	if !exists {
		// Создаем бакет, если он не существует
		err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: region,
		})
		if err != nil {
			return fmt.Errorf("ошибка при создании бакета: %w", err)
		}
	}
	return nil
}

// SaveImage сохраняет изображение в MinIO бакет и возвращает URL объекта
func (m *MinioClient) SaveImage(ctx context.Context, bucketName, objectName string, imageData []byte) (string, error) {
	// Проверяем существование бакета
	if err := m.ensureBucketExists(ctx, bucketName, m.conf.Region); err != nil {
		return "", models.NewImageStorageError(fmt.Sprintf("ошибка при проверке или создании бакета %s", bucketName), err)
	}

	// Загружаем изображение в бакет
	if _, err := m.client.PutObject(ctx, bucketName, objectName, bytes.NewReader(imageData), int64(len(imageData)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	}); err != nil {
		return "", models.NewImageStorageError(
			fmt.Sprintf("ошибка при загрузке изображения в MinIO в бакет %s с объектом %s", bucketName, objectName), err,
		)
	}

	// Формируем URL для объекта
	url := fmt.Sprintf("http://%s/%s/%s", m.client.EndpointURL().Host, bucketName, objectName)
	return url, nil
}
