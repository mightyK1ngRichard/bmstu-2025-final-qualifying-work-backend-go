package handler_test

import (
	handler "2025_CakeLand_API/internal/pkg/auth/delivery/grpc"
	"2025_CakeLand_API/internal/pkg/auth/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/auth/mocks"
	umodels "2025_CakeLand_API/internal/pkg/auth/usecase/models"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func TestRegisterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаём мок для IAuthUsecase
	mockAuthUsecase := mocks.NewMockIAuthUsecase(ctrl)

	// Создаём gRPC-хэндлер с мокнутым usecase
	h := handler.NewGrpcAuthHandler(mockAuthUsecase)

	// Настроим мок: если вызывается Register, он возвращает успешный результат
	mockAuthUsecase.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, req umodels.RegisterReq) {
			// Извлекаем метаданные из контекста для проверки
			md, _ := metadata.FromIncomingContext(ctx)
			assert.Equal(t, "some-fingerprint-value", md["fingerprint"][0])
		}).
		Return(&umodels.RegisterRes{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresIn:    time.Time{},
		}, nil)

	// Добавляем метаданные с fingerprint в контекст
	md := metadata.New(map[string]string{
		"fingerprint": "some-fingerprint-value",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	// Вызываем gRPC-хэндлер
	res, err := h.Register(ctx, &generated.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	// Проверяем результат
	assert.NoError(t, err)
	assert.Equal(t, "test-access-token", res.AccessToken)
	assert.Equal(t, "test-refresh-token", res.RefreshToken)
}
