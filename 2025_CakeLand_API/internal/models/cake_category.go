package models

import (
	"github.com/google/uuid"
	"time"
)

type CakeCategory struct {
	ID           uuid.UUID
	DateCreation time.Time
	CategoryID   uuid.UUID
	CakeID       uuid.UUID
}
