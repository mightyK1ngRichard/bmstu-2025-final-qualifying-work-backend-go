package metadata

import (
	"2025_CakeLand_API/internal/models"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"strings"
)

type MetadataKey string

const (
	KeyFingerprint   MetadataKey = "fingerprint"
	KeyAuthorization MetadataKey = "authorization"
)

type MetadataProvider struct {
}

func NewMetadataProvider() *MetadataProvider {
	return &MetadataProvider{}
}

func (m *MetadataProvider) GetValue(ctx context.Context, key MetadataKey) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("%w: %s", models.ErrNoMetadata, key)
	}

	val := md.Get(string(key))
	if len(val) == 0 {
		return "", fmt.Errorf("%w: %s", models.ErrNoMetadata, key)
	}

	// Если это Authorization заголовок, сохраняем токен без префикса
	if key == KeyAuthorization {
		return removeBearerPrefix(val[0]), nil
	}
	return val[0], nil
}

func (m *MetadataProvider) GetValues(ctx context.Context, keys ...MetadataKey) (map[MetadataKey]string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%w: metadata is missing", models.ErrNoMetadata)
	}

	values := make(map[MetadataKey]string)
	for _, key := range keys {
		val := md.Get(string(key))
		if len(val) == 0 {
			return nil, fmt.Errorf("%w: missing value for key %s", models.ErrNoMetadata, key)
		}

		// Если это Authorization заголовок, сохраняем токен без префикса
		if key == KeyAuthorization {
			values[key] = removeBearerPrefix(val[0])
		} else {
			values[key] = val[0]
		}
	}

	return values, nil
}

// removeBearerPrefix Проверяем, начинается ли строка с "Bearer" и убираем этот префикс
func removeBearerPrefix(token string) string {
	if strings.HasPrefix(token, "Bearer ") {
		return token[len("Bearer "):]
	}
	return token
}
