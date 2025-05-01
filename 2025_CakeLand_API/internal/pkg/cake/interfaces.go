package cake

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/cake/dto"
	ms "2025_CakeLand_API/internal/pkg/minio"
	"context"
	"github.com/google/uuid"
)

type ICakeUsecase interface {
	Cake(context.Context, dto.GetCakeReq) (*dto.GetCakeRes, error)
	CreateCake(context.Context, dto.CreateCakeReq) (*dto.CreateCakeRes, error)
	CreateFilling(context.Context, dto.CreateFillingReq) (*dto.CreateFillingRes, error)
	CreateCategory(context.Context, *dto.CreateCategoryReq) (*dto.CreateCategoryRes, error)
	Categories(context.Context) (*[]models.Category, error)
	Fillings(context.Context) (*[]models.Filling, error)
	AddCakeColor(context.Context, uuid.UUID, []string) error
	GetColors(context.Context) ([]string, error)
	GetCakesPreview(context.Context) ([]dto.PreviewCake, error)
	CategoryIDsByGenderName(context.Context, models.CategoryGender) ([]models.Category, error)
	CategoryPreviewCakes(context.Context, uuid.UUID) ([]*dto.PreviewCake, error)
}

type ICakeRepository interface {
	CakeByID(context.Context, dto.GetCakeReq) (*dto.GetCakeRes, error)
	CakeCategoriesIDs(context.Context, uuid.UUID) ([]uuid.UUID, error)
	CakeFillingsIDs(context.Context, uuid.UUID) ([]uuid.UUID, error)
	FillingByID(context.Context, uuid.UUID) (*models.Filling, error)
	CategoryByID(context.Context, uuid.UUID) (*models.Category, error)
	CakeImages(context.Context, uuid.UUID) ([]models.CakeImage, error)

	CreateCake(context.Context, dto.CreateCakeDBReq) error
	CreateFilling(context.Context, models.Filling) error
	CreateCategory(context.Context, *models.Category) error
	GetColors(context.Context) ([]string, error)
	AddCakeColor(context.Context, models.CakeColor) error

	GetCakeColorsByCakeID(context.Context, uuid.UUID) ([]models.CakeColor, error)
	GetUserByID(context.Context, uuid.UUID) (dto.Owner, error)
	GetCakesPreview(context.Context) ([]dto.PreviewCake, error)
	Categories(context.Context) (*[]models.Category, error)
	Fillings(context.Context) (*[]models.Filling, error)
	CategoryIDsByGenderName(context.Context, models.CategoryGender) ([]dto.DBCategory, error)
	CategoryCakesIDs(context.Context, uuid.UUID) ([]uuid.UUID, error)
	PreviewCakeByID(context.Context, uuid.UUID) (*dto.PreviewCake, error)
}

type IImageStorage interface {
	SaveImages(
		ctx context.Context,
		bucketName string,
		images map[ms.ImageID][]byte,
	) (map[ms.ImageID]string, error)
	SaveImage(
		ctx context.Context,
		bucketName string,
		objectName ms.ImageID,
		imageData []byte,
	) (string, error)
}
