package dto

import (
	genCake "2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/cake/dto"
	"2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
)

type UserInfo struct {
	User  Profile
	Cakes []dto.PreviewCake
}

func (u *UserInfo) ConvertToGrpcModel() *generated.UserInfo {
	cakes := make([]*genCake.PreviewCake, len(u.Cakes))
	for i, cake := range u.Cakes {
		cakes[i] = cake.ConvertToGrpcModel()
	}

	return &generated.UserInfo{
		User:  u.User.ConvertToGrpcModel(),
		Cakes: cakes,
	}
}
