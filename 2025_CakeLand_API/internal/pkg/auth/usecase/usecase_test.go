package usecase

import (
	"2025_CakeLand_API/internal/pkg/auth/dto"
	"2025_CakeLand_API/internal/pkg/auth/mocks"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"2025_CakeLand_API/internal/pkg/utils/logger"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthUsecase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := logger.NewLogger("local")
	mockRepo := mocks.NewMockIAuthRepository(ctrl)
	tokenator := jwt.NewTokenator()
	uc := NewAuthUsecase(log, tokenator, mockRepo)

	mockRepo.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(nil)

	t.Run("Success case", func(t *testing.T) {
		res, err := uc.Register(context.Background(), dto.RegisterReq{
			Email:       "test@example.com",
			Password:    "Password1",
			Fingerprint: "some-fingerprint",
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.AccessToken)
		assert.NotEmpty(t, res.RefreshToken)
		assert.NotEmpty(t, res.ExpiresIn)
	})
}
