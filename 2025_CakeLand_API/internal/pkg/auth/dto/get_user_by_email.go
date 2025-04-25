package dto

import (
	"2025_CakeLand_API/internal/models"
	"github.com/google/uuid"
)

type GetUserByEmailReq struct {
	Email string
}

type GetUserByEmailRes struct {
	ID               uuid.UUID
	Email            string
	PasswordHash     []byte
	RefreshTokensMap models.RefreshTokenMap // key: fingerprint, value: refreshToken
}
