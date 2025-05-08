package minio

import (
	"2025_CakeLand_API/internal/pkg/config"
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioProvider struct {
	client     *minio.Client
	conf       *config.MinioConfig
	bucketName string
}

type ImageID string

func NewMinioProvider(conf *config.MinioConfig) (*MinioProvider, error) {
	// Создаем нового клиента MinIO
	client, err := minio.New(conf.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания клиента MinIO: %w", err)
	}

	return &MinioProvider{
		client:     client,
		conf:       conf,
		bucketName: conf.Bucket,
	}, nil
}

func (m *MinioProvider) SaveImage(
	ctx context.Context,
	objectName ImageID,
	imageData []byte,
) (string, error) {
	// Проверяем существование бакета
	if err := m.ensureBucketExists(ctx, m.bucketName, m.conf.Region); err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("ошибка при проверке или создании бакета %s", m.bucketName))
	}

	// Формируем путь
	objectPath := string(objectName)

	// Загружаем изображение
	if _, err := m.client.PutObject(ctx, m.bucketName, objectPath, bytes.NewReader(imageData), int64(len(imageData)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	}); err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("ошибка при загрузке изображения в MinIO в бакет %s с объектом %s", m.bucketName, objectPath))
	}

	// Формируем URL
	url := fmt.Sprintf("http://%s/%s/%s", m.client.EndpointURL().Host, m.bucketName, objectPath)
	return url, nil
}

func (m *MinioProvider) SaveImages(
	ctx context.Context,
	images map[ImageID][]byte,
) (map[ImageID]string, error) {
	// Проверяем существование бакета
	if err := m.ensureBucketExists(ctx, m.bucketName, m.conf.Region); err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("ошибка при проверке или создании бакета %s", m.bucketName))
	}

	// Карта для хранения URL-ов загруженных изображений
	urls := make(map[ImageID]string)
	errs := make(chan error, len(images))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for objectName, imageData := range images {
		wg.Add(1)
		go func(objectName ImageID, imageData []byte) {
			defer wg.Done()

			// Формируем путь: objectName
			objectPath := string(objectName)

			// Загружаем изображение в бакет
			_, err := m.client.PutObject(ctx, m.bucketName, objectPath, bytes.NewReader(imageData), int64(len(imageData)), minio.PutObjectOptions{
				ContentType: "image/jpeg",
			})
			if err != nil {
				errs <- errors.Wrapf(err, fmt.Sprintf("ошибка при загрузке изображения в MinIO в бакет %s с объектом %s", m.bucketName, objectName))
				return
			}

			// Формируем URL и добавляем в карту
			url := fmt.Sprintf("http://%s/%s/%s", m.client.EndpointURL().Host, m.bucketName, objectName)
			mu.Lock()
			urls[objectName] = url
			mu.Unlock()
		}(objectName, imageData)
	}

	wg.Wait()
	close(errs)

	// Проверяем ошибки
	for err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return urls, nil
}

func (m *MinioProvider) SaveFile(
	ctx context.Context,
	objectName string,
	data []byte,
	contentType string,
) (string, error) {
	// Проверка/создание бакета
	if err := m.ensureBucketExists(ctx, m.bucketName, m.conf.Region); err != nil {
		return "", errors.Wrapf(err, "ошибка при проверке или создании бакета %s", m.bucketName)
	}

	// Загрузка файла
	if _, err := m.client.PutObject(ctx, m.bucketName, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	}); err != nil {
		return "", errors.Wrapf(err, "ошибка при загрузке файла в бакет %s с объектом %s", m.bucketName, objectName)
	}

	// Формирование URL
	url := fmt.Sprintf("http://%s/%s/%s", m.client.EndpointURL().Host, m.bucketName, objectName)
	return url, nil
}

// ensureBucketExists проверяет, существует ли бакет, и создает его, если нет
func (m *MinioProvider) ensureBucketExists(ctx context.Context, bucketName string, region string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("ошибка при проверке существования бакета: %w", err)
	}
	if !exists {
		// Создаем бакет, если он не существует
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: region,
		})
		if err != nil {
			return fmt.Errorf("ошибка при создании бакета: %w", err)
		}
	}
	return nil
}
