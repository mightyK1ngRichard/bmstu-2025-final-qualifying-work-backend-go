package models

import (
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"github.com/google/uuid"
)

// Filling Модель начинки
type Filling struct {
	ID          uuid.UUID // Код
	Name        string    // Название
	ImageURL    string    // Картинка
	Content     string    // Содержимое начинки
	KgPrice     float64   // Цена за кг
	Description string    // Описание
}

func (f *Filling) ConvertToFillingGRPC() *generated.Filling {
	return &generated.Filling{
		Id:          f.ID.String(),
		Name:        f.Name,
		ImageUrl:    f.ImageURL,
		Content:     f.Content,
		KgPrice:     f.KgPrice,
		Description: f.Description,
	}
}
