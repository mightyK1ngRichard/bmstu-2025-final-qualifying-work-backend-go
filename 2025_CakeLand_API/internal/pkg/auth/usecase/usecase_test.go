package usecase

import (
	"2025_CakeLand_API/internal/pkg/auth/mocks"
	umodels "2025_CakeLand_API/internal/pkg/auth/usecase/models"
	"2025_CakeLand_API/internal/pkg/utils"
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
	validator := utils.NewValidator()
	uc := NewAuthUsecase(log, validator, mockRepo)

	mockRepo.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(nil)

	t.Run("Success case", func(t *testing.T) {
		res, err := uc.Register(context.Background(), umodels.RegisterReq{
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

	t.Run("Bad Email", func(t *testing.T) {
		res, err := uc.Register(context.Background(), umodels.RegisterReq{
			Email:       "testexample.com",
			Password:    "Password1!",
			Fingerprint: "some-fingerprint",
		})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Bad Password", func(t *testing.T) {
		res, err := uc.Register(context.Background(), umodels.RegisterReq{
			Email:       "test@example.com",
			Password:    "Passwor",
			Fingerprint: "some-fingerprint",
		})

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
