package handler_test

import (
	handler "2025_CakeLand_API/internal/pkg/auth/delivery/grpc"
	"2025_CakeLand_API/internal/pkg/auth/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/auth/dto"
	"2025_CakeLand_API/internal/pkg/auth/mocks"
	"2025_CakeLand_API/internal/pkg/utils"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestRegisterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаём мок для IAuthUsecase
	mockAuthUsecase := mocks.NewMockIAuthUsecase(ctrl)

	// Создаём gRPC-хэндлер с мокнутым usecase
	validator := utils.NewValidator()
	mdProvider := md.NewMetadataProvider()
	h := handler.NewGrpcAuthHandler(validator, mockAuthUsecase, mdProvider)

	// Настроим мок: если вызывается Register, он возвращает успешный результат
	mockAuthUsecase.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(&dto.RegisterRes{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresIn:    time.Time{},
		}, nil)

	// Добавляем метаданные с fingerprint в контекст
	md := metadata.New(map[string]string{
		"fingerprint": "some-fingerprint-value",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	t.Run("Success registration", func(t *testing.T) {
		// Вызываем gRPC-хэндлер
		res, err := h.Register(ctx, &generated.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		})

		// Проверяем результат
		assert.NoError(t, err)
		assert.Equal(t, "test-access-token", res.AccessToken)
		assert.Equal(t, "test-refresh-token", res.RefreshToken)
	})

	t.Run("Bad Email", func(t *testing.T) {
		res, err := h.Register(ctx, &generated.RegisterRequest{
			Email:    "test@example.com",
			Password: "Password1!",
		})

		assert.Error(t, err)
		assert.Nil(t, res)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("Bad Password", func(t *testing.T) {
		res, err := h.Register(ctx, &generated.RegisterRequest{
			Email:    "test@example.com",
			Password: "Passwor",
		})

		assert.Error(t, err)
		assert.Nil(t, res)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}
