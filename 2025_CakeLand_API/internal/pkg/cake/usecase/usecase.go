package usecase

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/cake"
	en "2025_CakeLand_API/internal/pkg/cake/entities"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"2025_CakeLand_API/internal/pkg/utils/sl"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"

	"github.com/google/uuid"
)

type CakeUseсase struct {
	log        *slog.Logger
	tokenator  *jwt.Tokenator
	repo       cake.ICakeRepository
	imageStore cake.IImageStorage
	bucketName string
}

func NewCakeUsecase(
	log *slog.Logger,
	tokenator *jwt.Tokenator,
	repo cake.ICakeRepository,
	imageStore cake.IImageStorage,
	bucketName string,
) *CakeUseсase {
	return &CakeUseсase{
		log:        log,
		tokenator:  tokenator,
		repo:       repo,
		imageStore: imageStore,
		bucketName: bucketName,
	}
}

func (u *CakeUseсase) Cake(ctx context.Context, in en.GetCakeReq) (*en.GetCakeRes, error) {
	res, err := u.repo.CakeByID(ctx, in)
	if err != nil {
		u.log.Error("[Usecase.Cake] ошибка получения торта по id из бд",
			slog.String("cakeID", in.CakeID.String()),
			sl.Err(err),
		)
		return nil, err
	}

	return &en.GetCakeRes{
		Cake: res.Cake,
	}, nil
}

func (u *CakeUseсase) CreateCake(ctx context.Context, in en.CreateCakeReq) (*en.CreateCakeRes, error) {
	// Достаём userID из токена если он не протух
	userID, err := u.tokenator.GetUserIDFromToken(in.AccessToken, false)
	if err != nil {
		u.log.Error(`[Usecase.CreateCake] ошибка получения userID из refresh токена`, sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	cakeID := uuid.New()
	// Добавляем изображение в хранилище
	imageURL, err := u.imageStore.SaveImage(ctx, u.bucketName, cakeID.String(), in.ImageData)
	if err != nil {
		u.log.Error(`[Usecase.CreateCake] ошибка загрузки изображения в хранилище`, sl.Err(err))
		return nil, err
	}

	// Создаём торт в бд
	if err = u.repo.CreateCake(ctx, in.ConvertToCreateCakeDBReq(cakeID.String(), userID, imageURL)); err != nil {
		u.log.Error(`[Usecase.CreateCake] ошибка сохранения торта в бд`, sl.Err(err))
		return nil, err
	}

	return &en.CreateCakeRes{
		CakeID: cakeID.String(),
	}, nil
}

func (u *CakeUseсase) CreateFilling(ctx context.Context, in en.CreateFillingReq) (*en.CreateFillingRes, error) {
	fillingID := uuid.New()
	// Добавляем изображение в хранилище
	imageURL, err := u.imageStore.SaveImage(ctx, u.bucketName, fillingID.String(), in.ImageData)
	if err != nil {
		u.log.Error("[Usecase.CreateFilling] ошибка загрузки изображения в хранилище", sl.Err(err))
		return nil, models.ErrInternal
	}

	filling := models.Filling{
		ID:          fillingID,
		Name:        in.Name,
		ImageURL:    imageURL,
		Content:     in.Content,
		KgPrice:     in.KgPrice,
		Description: in.Description,
	}
	err = u.repo.CreateFilling(ctx, filling)
	if err != nil {
		return nil, err
	}

	return &en.CreateFillingRes{
		Filling: filling,
	}, nil
}

func (u *CakeUseсase) CreateCategory(ctx context.Context, in *en.CreateCategoryReq) (*en.CreateCategoryRes, error) {
	categoryUUID := uuid.New()
	imageURL, err := u.imageStore.SaveImage(ctx, u.bucketName, categoryUUID.String(), in.ImageData)
	if err != nil {
		u.log.Error("[Usecase.CreateCategory] ошибка загрузки изображения в хранилище", sl.Err(err))
		return nil, models.ErrInternal
	}

	newCategory := models.Category{
		ID:       categoryUUID,
		Name:     in.Name,
		ImageURL: imageURL,
	}
	if err := u.repo.CreateCategory(ctx, &newCategory); err != nil {
		return nil, err
	}

	return &en.CreateCategoryRes{
		Category: newCategory,
	}, nil
}

func (u *CakeUseсase) Categories(ctx context.Context) (*[]models.Category, error) {
	return u.repo.Categories(ctx)
}

func (u *CakeUseсase) Fillings(ctx context.Context) (*[]models.Filling, error) {
	return u.repo.Fillings(ctx)
}

func (u *CakeUseсase) Cakes(ctx context.Context) (*[]models.Cake, error) {
	return u.repo.Cakes(ctx)
}
