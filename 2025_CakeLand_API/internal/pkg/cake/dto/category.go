package dto

import (
	"2025_CakeLand_API/internal/models"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type DBCategory struct {
	ID              uuid.UUID
	Name            null.String
	ImageURL        null.String
	CategoryGenders []models.CategoryGender
}

func (c *DBCategory) ConvertToCategory() models.Category {
	return models.Category{
		ID:              c.ID,
		Name:            c.Name.String,
		ImageURL:        c.ImageURL.String,
		CategoryGenders: c.CategoryGenders,
	}
}
