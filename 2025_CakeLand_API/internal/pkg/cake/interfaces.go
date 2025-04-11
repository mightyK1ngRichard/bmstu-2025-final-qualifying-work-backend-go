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
	Cakes(context.Context) (*[]models.Cake, error)
	CategoryIDsByGenderName(context.Context, models.CategoryGender) ([]models.Category, error)
	CategoryPreviewCakes(context.Context, uuid.UUID) ([]*dto.PreviewCake, error)
}

type ICakeRepository interface {
	CakeByID(context.Context, dto.GetCakeReq) (*dto.GetCakeRes, error)
	CreateCake(context.Context, dto.CreateCakeDBReq) error
	CreateFilling(context.Context, models.Filling) error
	CreateCategory(context.Context, *models.Category) error
	Categories(context.Context) (*[]models.Category, error)
	Fillings(context.Context) (*[]models.Filling, error)
	Cakes(context.Context) (*[]models.Cake, error)
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
