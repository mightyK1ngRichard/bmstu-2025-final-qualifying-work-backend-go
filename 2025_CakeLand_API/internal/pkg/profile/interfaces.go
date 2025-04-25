package profile

import (
	"2025_CakeLand_API/internal/models"
	cakeDto "2025_CakeLand_API/internal/pkg/cake/dto"
	"2025_CakeLand_API/internal/pkg/profile/dto"
	"context"
	"github.com/google/uuid"
)

type IProfileUsecase interface {
	UserInfo(context.Context, string) (*dto.UserInfo, error)
	UserInfoByID(context.Context, uuid.UUID) (*models.UserInfo, error)
}

type IProfileRepository interface {
	UserInfo(context.Context, uuid.UUID) (*dto.Profile, error)
	CakesByUserID(ctx context.Context, userID uuid.UUID) ([]cakeDto.PreviewCakeDB, error)
}
