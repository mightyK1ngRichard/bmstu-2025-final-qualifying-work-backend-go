package models

import (
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"github.com/google/uuid"
)

// Category Модель категории
type Category struct {
	ID       uuid.UUID // Код
	Name     string    // Название
	ImageURL string    // Картинка
}

func (c *Category) ConvertToCategoryGRPC() *generated.Category {
	return &generated.Category{
		Id:       c.ID.String(),
		Name:     c.Name,
		ImageUrl: c.ImageURL,
	}
}
