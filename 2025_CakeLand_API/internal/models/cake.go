package models

import (
	"2025_CakeLand_API/internal/models/errs"
	gen "2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type CakeStatus string

const (
	CakeStatusPending  CakeStatus = "pending"
	CakeStatusApproved CakeStatus = "approved"
	CakeStatusRejected CakeStatus = "rejected"
	CakeStatusHidden   CakeStatus = "hidden"
)

// Cake Модель торта
type Cake struct {
	ID              uuid.UUID   // Код
	Name            string      // Название
	PreviewImageURL string      // Картинка товара
	KgPrice         float64     // Цена за кг
	ReviewsCount    int32       // Количество отзывов
	StarsSum        int32       // Сумма звёзд
	Description     string      // Описание
	Mass            float64     // Масса торта
	Status          CakeStatus  // Статус торта
	Model3DURL      null.String // Ссылка на 3D модель
	DateCreation    time.Time   // Дата создания торта
	DiscountKgPrice null.Float  // Скидочная цена за кг
	DiscountEndTime null.Time   // Дата окончания скидки
	Owner           User        // Владелец
	Fillings        []Filling   // Слои торта
	Categories      []Category  // Категории торта
	Images          []CakeImage // Фотографии торта
	CakeColor       []CakeColor // Цвета торта
}

type CakeColor struct {
	ID        uuid.UUID
	HexString string
	CakeID    uuid.UUID
}

type CakeImage struct {
	ID       uuid.UUID
	ImageURL null.String
}

// Model -> Proto

func (s CakeStatus) ToProto() gen.CakeStatus {
	switch s {
	case CakeStatusPending:
		return gen.CakeStatus_PENDING
	case CakeStatusApproved:
		return gen.CakeStatus_APPROVED
	case CakeStatusRejected:
		return gen.CakeStatus_REJECTED
	case CakeStatusHidden:
		return gen.CakeStatus_HIDDEN
	default:
		return gen.CakeStatus_CAKE_STATUS_UNSPECIFIED
	}
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

	// рейтинг = сумма / кол-во
	var rating int32 = 0
	if c.ReviewsCount != 0 {
		rating = c.StarsSum / c.ReviewsCount
	}

	var model3DURL *string
	if c.Model3DURL.Valid {
		model3DURL = &c.Model3DURL.String
	}

	return &gen.Cake{
		Id:              c.ID.String(),
		Name:            c.Name,
		ImageUrl:        c.PreviewImageURL,
		KgPrice:         c.KgPrice,
		Rating:          rating,
		Description:     c.Description,
		Mass:            c.Mass,
		Status:          c.Status.ToProto(),
		Owner:           c.Owner.ConvertToUserGRPC(),
		Fillings:        grpcFillings,
		Categories:      grpcCategories,
		DiscountKgPrice: discountKgPrice,
		DiscountEndTime: discountEndTime,
		DateCreation:    timestamppb.New(c.DateCreation),
		Images:          cakeImages,
		ReviewsCount:    c.ReviewsCount,
		Model3DURL:      model3DURL,
	}
}

// Proto -> Model

func FromProtoCakeStatus(status gen.CakeStatus) (CakeStatus, error) {
	switch status {
	case gen.CakeStatus_CAKE_STATUS_UNSPECIFIED:
		return "", errs.ErrUnknownCakeStatus
	case gen.CakeStatus_APPROVED:
		return CakeStatusApproved, nil
	case gen.CakeStatus_REJECTED:
		return CakeStatusRejected, nil
	case gen.CakeStatus_HIDDEN:
		return CakeStatusHidden, nil
	case gen.CakeStatus_PENDING:
		return CakeStatusPending, nil
	default:
		return "", errs.ErrUnknownCakeStatus
	}
}

// SQL

func (s *CakeStatus) Scan(value interface{}) error {
	var str string

	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("cannot scan type %T into CakeStatus", value)
	}

	switch CakeStatus(str) {
	case CakeStatusPending, CakeStatusApproved, CakeStatusRejected, CakeStatusHidden:
		*s = CakeStatus(str)
		return nil
	default:
		return fmt.Errorf("invalid CakeStatus: %s", str)
	}
}

func (s CakeStatus) Value() (driver.Value, error) {
	switch s {
	case CakeStatusPending, CakeStatusApproved, CakeStatusRejected, CakeStatusHidden:
		return string(s), nil
	default:
		return nil, fmt.Errorf("invalid CakeStatus: %s", s)
	}
}
