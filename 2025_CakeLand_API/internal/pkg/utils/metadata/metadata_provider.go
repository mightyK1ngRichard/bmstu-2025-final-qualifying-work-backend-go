package metadata

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models/errs"
	"context"
	"google.golang.org/grpc/metadata"
	"strings"
)

type MetadataProvider struct {
}

func NewMetadataProvider() *MetadataProvider {
	return &MetadataProvider{}
}

func (m *MetadataProvider) GetValue(ctx context.Context, key domains.MetadataKey) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errs.ErrNoMetadata
	}

	val := md.Get(key.String())
	if len(val) == 0 {
		return "", errs.ErrNoMetadata
	}

	// Если это Authorization заголовок, сохраняем токен без префикса
	if key == domains.KeyAuthorization {
		return removeBearerPrefix(val[0]), nil
	}

	return val[0], nil
}

func (m *MetadataProvider) GetValues(ctx context.Context, keys ...domains.MetadataKey) (map[domains.MetadataKey]string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errs.ErrNoMetadata
	}

	values := make(map[domains.MetadataKey]string)
	for _, key := range keys {
		val := md.Get(key.String())
		if len(val) == 0 {
			return nil, errs.ErrNoMetadata
		}

		// Если это Authorization заголовок, сохраняем токен без префикса
		if key == domains.KeyAuthorization {
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
