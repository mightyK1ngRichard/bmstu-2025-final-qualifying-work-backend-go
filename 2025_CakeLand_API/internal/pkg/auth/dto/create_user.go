package dto

import (
	"2025_CakeLand_API/internal/models"
	"github.com/google/uuid"
)

type CreateUserReq struct {
	UUID             uuid.UUID
	Email            string
	Nickname         string
	PasswordHash     []byte
	RefreshTokensMap models.RefreshTokenMap
}
