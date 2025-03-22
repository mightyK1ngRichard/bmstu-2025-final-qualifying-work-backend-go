package entities

import (
	"2025_CakeLand_API/internal/models"
	gen "2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	ms "2025_CakeLand_API/internal/pkg/s3storage"

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
	Name             string   // Название торта
	PreviewImageData []byte   // Данные изображения торта
	KgPrice          float64  // Цена за кг
	Rating           int32    // Рейтинг (0-5)
	Description      string   // Описание торта
	Mass             float64  // Масса торта
	IsOpenForSale    bool     // Доступен ли для продажи
	FillingIDs       []string // Список ID начинок
	CategoryIDs      []string // Список ID категорий
	AccessToken      string   // Токен пользователя
	Images           [][]byte // Изображения торта
}

func NewCreateCakeReq(in *gen.CreateCakeRequest, accessToken string) CreateCakeReq {
	return CreateCakeReq{
		Name:             in.Name,
		PreviewImageData: in.PreviewImageData,
		KgPrice:          in.KgPrice,
		Rating:           in.Rating,
		Description:      in.Description,
		Mass:             in.Mass,
		IsOpenForSale:    in.IsOpenForSale,
		FillingIDs:       in.FillingIds,
		CategoryIDs:      in.CategoryIds,
		AccessToken:      accessToken,
		Images:           in.Images,
	}
}

func (req *CreateCakeReq) ConvertToCreateCakeDBReq(
	cakeID string,
	previewImageURL string,
	ownerID string,
	images map[ms.ImageID]string,
) CreateCakeDBReq {
	return CreateCakeDBReq{
		ID:              cakeID,
		Name:            req.Name,
		PreviewImageURL: previewImageURL,
		KgPrice:         req.KgPrice,
		Rating:          req.Rating,
		Description:     req.Description,
		Mass:            req.Mass,
		IsOpenForSale:   req.IsOpenForSale,
		OwnerID:         ownerID,
		FillingIDs:      req.FillingIDs,
		CategoryIDs:     req.CategoryIDs,
		Images:          images,
	}
}

type CreateCakeDBReq struct {
	ID              string                // Код торта
	Name            string                // Название торта
	PreviewImageURL string                // URL изображения торта
	KgPrice         float64               // Цена за кг
	Rating          int32                 // Рейтинг (0-5)
	Description     string                // Описание торта
	Mass            float64               // Масса торта
	IsOpenForSale   bool                  // Доступен ли для продажи
	OwnerID         string                // ID владельца
	FillingIDs      []string              // Список ID начинок
	CategoryIDs     []string              // Список ID категорий
	Images          map[ms.ImageID]string // Фотографии торта
}

type CreateCakeRes struct {
	CakeID string
}

// CreateFilling

type CreateFillingReq struct {
	Name        string  // Название начинки
	ImageData   []byte  // Картинка начинки
	Content     string  // Содержимое начинки
	KgPrice     float64 // Цена за кг начинки
	Description string  // Описание начинки
	AccessToken string  // Токен пользователя
}

type CreateFillingRes struct {
	Filling models.Filling
}

// Create Category

type CreateCategoryReq struct {
	Name        string // Название категории
	ImageData   []byte // Фотография категории
	AccessToken string // Токен пользователя
}

type CreateCategoryRes struct {
	Category models.Category
}
