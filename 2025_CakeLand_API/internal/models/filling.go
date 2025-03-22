package models

import (
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"github.com/google/uuid"
	"github.com/guregu/null"
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

type DBFilling struct {
	ID          null.String
	Name        null.String
	ImageURL    null.String
	Content     null.String
	KgPrice     null.Float
	Description null.String
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

func (f *DBFilling) ConvertToFilling() *Filling {
	if !f.ID.Valid {
		return nil
	}

	return &Filling{
		ID:          uuid.MustParse(f.ID.String),
		Name:        f.Name.String,
		ImageURL:    f.ImageURL.String,
		Content:     f.Content.String,
		KgPrice:     f.KgPrice.Float64,
		Description: f.Description.String,
	}
}
