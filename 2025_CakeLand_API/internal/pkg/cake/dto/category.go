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
	CategoryGenders []string
}

func (c *DBCategory) ConvertToCategory() models.Category {
	var categoryGenders []models.CategoryGender
	for _, gender := range c.CategoryGenders {
		categoryGenders = append(categoryGenders, models.ConvertToCategoryGender(gender))
	}

	return models.Category{
		ID:              c.ID,
		Name:            c.Name.String,
		ImageURL:        c.ImageURL.String,
		CategoryGenders: categoryGenders,
	}
}
