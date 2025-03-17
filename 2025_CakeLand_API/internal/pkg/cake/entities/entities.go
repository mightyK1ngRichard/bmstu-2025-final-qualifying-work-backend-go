package entities

import (
	"2025_CakeLand_API/internal/models"

	"github.com/google/uuid"
)

// GetCake

type GetCakeReq struct {
	CakeID uuid.UUID
}

type GetCakeRes struct {
	Cake models.Cake
}

// CreateCake

type CreateCakeReq struct {
	Name          string   // Название торта
	ImageData     []byte   // Данные изображения торта
	KgPrice       float64  // Цена за кг
	Rating        int32    // Рейтинг (0-5)
	Description   string   // Описание торта
	Mass          float64  // Масса торта
	IsOpenForSale bool     // Доступен ли для продажи
	FillingIDs    []string // Список ID начинок
	CategoryIDs   []string // Список ID категорий
	AccessToken   string   // Токен пользователя
}

func (req *CreateCakeReq) ConvertToCreateCakeDBReq(
	cakeID string,
	ownerID string,
	imageURL string,
) CreateCakeDBReq {
	return CreateCakeDBReq{
		ID:            cakeID,
		Name:          req.Name,
		ImageURL:      imageURL,
		KgPrice:       req.KgPrice,
		Rating:        req.Rating,
		Description:   req.Description,
		Mass:          req.Mass,
		IsOpenForSale: req.IsOpenForSale,
		OwnerID:       ownerID,
		FillingIDs:    req.FillingIDs,
		CategoryIDs:   req.CategoryIDs,
	}
}

type CreateCakeDBReq struct {
	ID            string   // Код торта
	Name          string   // Название торта
	ImageURL      string   // URL изображения торта
	KgPrice       float64  // Цена за кг
	Rating        int32    // Рейтинг (0-5)
	Description   string   // Описание торта
	Mass          float64  // Масса торта
	IsOpenForSale bool     // Доступен ли для продажи
	OwnerID       string   // ID владельца
	FillingIDs    []string // Список ID начинок
	CategoryIDs   []string // Список ID категорий
}

type CreateCakeRes struct {
	CakeID string
}

// CreateFilling

type CreateFillingReq struct {
	Name        string  // Название
	ImageData   []byte  // Картинка
	Content     string  // Содержимое начинки
	KgPrice     float64 // Цена за кг
	Description string  // Описание
}

type CreateFillingRes struct {
	Filling models.Filling
}

// Create Category

type CreateCategoryReq struct {
	Name      string
	ImageData []byte
}

type CreateCategoryRes struct {
	Category models.Category
}
