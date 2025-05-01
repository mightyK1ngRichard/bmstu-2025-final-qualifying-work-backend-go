package profile

import (
	"2025_CakeLand_API/internal/models"
	cakeDto "2025_CakeLand_API/internal/pkg/cake/dto"
	gen "2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/profile/dto"
	"context"
	"github.com/google/uuid"
)

type IProfileUsecase interface {
	UserInfo(context.Context, string) (*dto.UserInfo, error)
	UserInfoByID(context.Context, uuid.UUID) (*models.UserInfo, error)
	CreateAddress(context.Context, string, *models.Address) (*models.Address, error)
	GetUserAddresses(context.Context, string) ([]models.Address, error)
	UpdateUserAddresses(context.Context, string, *gen.UpdateUserAddressesReq) (models.Address, error)
}

type IProfileRepository interface {
	UserInfo(context.Context, uuid.UUID) (*dto.Profile, error)
	CakesByUserID(ctx context.Context, userID uuid.UUID) ([]cakeDto.PreviewCakeDB, error)
	CreateAddress(context.Context, *models.Address) error
	GetUserAddresses(context.Context, uuid.UUID) ([]models.Address, error)
	UpdateUserAddresses(context.Context, uuid.UUID, *gen.UpdateUserAddressesReq) (models.Address, error)
}
