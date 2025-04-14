package models

import (
	gen "2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// Cake Модель торта
type Cake struct {
	ID              uuid.UUID   // Код
	Name            string      // Название
	PreviewImageURL string      // Картинка товара
	KgPrice         float64     // Цена за кг
	Rating          int         // Рейтинг (от 0 до 5)
	Description     string      // Описание
	Mass            float64     // Масса торта
	IsOpenForSale   bool        // Флаг возможности продажи торта
	DateCreation    time.Time   // Дата создания торта
	DiscountKgPrice null.Float  // Скидочная цена за кг
	DiscountEndTime null.Time   // Дата окончания скидки
	Owner           User        // Владелец
	Fillings        []Filling   // Слои торта
	Categories      []Category  // Категории торта
	Images          []CakeImage // Фотографии торта
}

type CakeImage struct {
	ID       uuid.UUID
	ImageURL null.String
}

func (c *CakeImage) ConvertToCakeImageGRPC() *gen.Cake_CakeImage {
	return &gen.Cake_CakeImage{
		Id:       c.ID.String(),
		ImageUrl: c.ImageURL.String,
	}
}

func (c *Cake) ConvertToCakeGRPC() *gen.Cake {
	grpcFillings := make([]*gen.Filling, len(c.Fillings))
	for i, it := range c.Fillings {
		grpcFillings[i] = it.ConvertToFillingGRPC()
	}
	grpcCategories := make([]*gen.Category, len(c.Categories))
	for i, it := range c.Categories {
		grpcCategories[i] = it.ConvertToCategoryGRPC()
	}

	var discountKgPrice *float64
	if c.DiscountKgPrice.Valid {
		val := c.DiscountKgPrice.Float64
		discountKgPrice = &val
	}

	var discountEndTime *timestamppb.Timestamp
	if c.DiscountEndTime.Valid {
		discountEndTime = timestamppb.New(c.DiscountEndTime.Time)
	}

	cakeImages := make([]*gen.Cake_CakeImage, len(c.Images))
	for i, it := range c.Images {
		cakeImages[i] = it.ConvertToCakeImageGRPC()
	}

	return &gen.Cake{
		Id:              c.ID.String(),
		Name:            c.Name,
		ImageUrl:        c.PreviewImageURL,
		KgPrice:         c.KgPrice,
		Rating:          int32(c.Rating),
		Description:     c.Description,
		Mass:            c.Mass,
		IsOpenForSale:   c.IsOpenForSale,
		Owner:           c.Owner.ConvertToUserGRPC(),
		Fillings:        grpcFillings,
		Categories:      grpcCategories,
		DiscountKgPrice: discountKgPrice,
		DiscountEndTime: discountEndTime,
		DateCreation:    timestamppb.New(c.DateCreation),
		Images:          cakeImages,
	}
}
