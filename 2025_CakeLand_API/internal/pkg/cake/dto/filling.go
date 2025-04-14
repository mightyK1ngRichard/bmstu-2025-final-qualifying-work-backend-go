package dto

import (
	"2025_CakeLand_API/internal/models"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type DBFilling struct {
	ID          uuid.UUID
	Name        null.String
	ImageURL    null.String
	Content     null.String
	KgPrice     null.Float
	Description null.String
}

func (f *DBFilling) ConvertToFilling() models.Filling {
	return models.Filling{
		ID:          f.ID,
		Name:        f.Name.String,
		ImageURL:    f.ImageURL.String,
		Content:     f.Content.String,
		KgPrice:     f.KgPrice.Float64,
		Description: f.Description.String,
	}
}
