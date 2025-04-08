package usecase

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/cake"
	en "2025_CakeLand_API/internal/pkg/cake/dto"
	ms "2025_CakeLand_API/internal/pkg/minio"
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

	// Добавляем изображение в хранилище
	images := make(map[ms.ImageID][]byte, len(in.Images)+1) // Капасити это фотографии тортов + превью фотография
	for _, imageData := range in.Images {
		cakeID := ms.ImageID(uuid.New().String())
		images[cakeID] = imageData
	}
	previewImageID := ms.ImageID(uuid.New().String())
	images[previewImageID] = in.PreviewImageData
	res, err := u.imageStore.SaveImages(ctx, u.bucketName, images)

	// Получаем preview
	previewImageURL, ok := res[previewImageID]
	if !ok {
		u.log.Error("[Usecase.CreateFilling] ошибка получения preview по ключу. Текст ошибки: ", sl.Err(err))
		return nil, models.ErrPreviewImageNotFound
	}

	// Создаём торт в бд
	cakeID := uuid.New()
	if err = u.repo.CreateCake(ctx, in.ConvertToCreateCakeDBReq(cakeID.String(), previewImageURL, userID, res)); err != nil {
		u.log.Error(`[Usecase.CreateCake] ошибка сохранения торта в бд`, sl.Err(err))
		return nil, err
	}

	return &en.CreateCakeRes{
		CakeID: cakeID.String(),
	}, nil
}

func (u *CakeUseсase) CreateFilling(ctx context.Context, in en.CreateFillingReq) (*en.CreateFillingRes, error) {
	// Достаём userID из токена если он не протух
	_, err := u.tokenator.GetUserIDFromToken(in.AccessToken, false)
	if err != nil {
		u.log.Error(`[Usecase.CreateCake] ошибка получения userID из refresh токена`, sl.Err(err))
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	fillingID := uuid.New()
	// Добавляем изображение в хранилище
	imageURL, err := u.imageStore.SaveImage(ctx, u.bucketName, ms.ImageID(fillingID.String()), in.ImageData)
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
	// Достаём userID из токена если он не протух
	_, err := u.tokenator.GetUserIDFromToken(in.AccessToken, false)
	if err != nil {
		u.log.Error(`[Usecase.CreateCake] ошибка получения userID из refresh токена`, sl.Err(err))
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	categoryUUID := uuid.New()
	imageURL, err := u.imageStore.SaveImage(ctx, u.bucketName, ms.ImageID(categoryUUID.String()), in.ImageData)
	if err != nil {
		u.log.Error("[Usecase.CreateCategory] ошибка загрузки изображения в хранилище", sl.Err(err))
		return nil, models.ErrInternal
	}

	newCategory := models.Category{
		ID:       categoryUUID,
		Name:     in.Name,
		ImageURL: imageURL,
	}
	if err = u.repo.CreateCategory(ctx, &newCategory); err != nil {
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

func (u *CakeUseсase) CategoryIDsByGenderName(ctx context.Context, genTag models.CategoryGender) ([]models.Category, error) {
	dbCategories, err := u.repo.CategoryIDsByGenderName(ctx, genTag)
	if err != nil {
		return nil, err
	}

	categories := make([]models.Category, len(dbCategories))
	for i, category := range dbCategories {
		categories[i] = category.ConvertToCategory()
	}
	return categories, nil
}
