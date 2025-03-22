package models

import (
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

// Category Модель категории
type Category struct {
	ID       uuid.UUID // Код
	Name     string    // Название
	ImageURL string    // Картинка
}

type DBCategory struct {
	ID       null.String
	Name     null.String
	ImageURL null.String
}

func (c *Category) ConvertToCategoryGRPC() *generated.Category {
	return &generated.Category{
		Id:       c.ID.String(),
		Name:     c.Name,
		ImageUrl: c.ImageURL,
	}
}

func (c *DBCategory) ConvertToCategory() *Category {
	id, err := uuid.Parse(c.ID.String)
	if err != nil {
		return nil
	}

	return &Category{
		ID:       id,
		Name:     c.Name.String,
		ImageURL: c.ImageURL.String,
	}
}
