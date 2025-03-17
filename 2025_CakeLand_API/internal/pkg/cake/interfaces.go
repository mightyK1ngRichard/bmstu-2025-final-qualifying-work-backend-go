package cake

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/cake/entities"
	"context"
)

type ICakeUsecase interface {
	Cake(context.Context, entities.GetCakeReq) (*entities.GetCakeRes, error)
	CreateCake(context.Context, entities.CreateCakeReq) (*entities.CreateCakeRes, error)
	CreateFilling(context.Context, entities.CreateFillingReq) (*entities.CreateFillingRes, error)
	CreateCategory(context.Context, *entities.CreateCategoryReq) (*entities.CreateCategoryRes, error)
	Categories(context.Context) (*[]models.Category, error)
	Fillings(context.Context) (*[]models.Filling, error)
	Cakes(context.Context) (*[]models.Cake, error)
}

type ICakeRepository interface {
	CakeByID(context.Context, entities.GetCakeReq) (*entities.GetCakeRes, error)
	CreateCake(context.Context, entities.CreateCakeDBReq) error
	CreateFilling(context.Context, models.Filling) error
	CreateCategory(context.Context, *models.Category) error
	Categories(context.Context) (*[]models.Category, error)
	Fillings(context.Context) (*[]models.Filling, error)
	Cakes(context.Context) (*[]models.Cake, error)
}

type IImageStorage interface {
	SaveImage(context.Context, string, string, []byte) (string, error)
}
