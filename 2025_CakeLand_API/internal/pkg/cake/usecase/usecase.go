package usecase

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/cake"
	"2025_CakeLand_API/internal/pkg/cake/dto"
	ms "2025_CakeLand_API/internal/pkg/minio"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"context"
	"github.com/google/uuid"
	"sync"
)

type CakeUseсase struct {
	tokenator  *jwt.Tokenator
	repo       cake.ICakeRepository
	imageStore cake.IImageStorage
	bucketName string
}

func NewCakeUsecase(
	tokenator *jwt.Tokenator,
	repo cake.ICakeRepository,
	imageStore cake.IImageStorage,
	bucketName string,
) *CakeUseсase {
	return &CakeUseсase{
		tokenator:  tokenator,
		repo:       repo,
		imageStore: imageStore,
		bucketName: bucketName,
	}
}

func (u *CakeUseсase) Cake(ctx context.Context, in dto.GetCakeReq) (*dto.GetCakeRes, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	res, err := u.repo.CakeByID(ctx, in)
	if err != nil {
		return nil, err
	}
	cakeInfo := res.Cake

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	errChan := make(chan error, 1)

	wg.Add(3)

	// Получаем категории торта
	go func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		ids, catIdsErr := u.repo.CakeCategoriesIDs(ctx, cakeInfo.ID)
		if catIdsErr != nil {
			trySendError(catIdsErr, errChan, cancel)
			return
		}

		// Получаем информацию каждой категории
		catMu := sync.Mutex{}
		catWg := sync.WaitGroup{}

		var cakeCategories []models.Category
		for _, categoryID := range ids {
			catWg.Add(1)
			go func() {
				defer catWg.Done()
				if ctx.Err() != nil {
					return
				}

				category, catErr := u.repo.CategoryByID(ctx, categoryID)
				if catErr != nil {
					trySendError(catErr, errChan, cancel)
					return
				}

				catMu.Lock()
				cakeCategories = append(cakeCategories, *category)
				catMu.Unlock()
			}()
		}

		catWg.Wait()
		if ctx.Err() != nil {
			return
		}

		mu.Lock()
		cakeInfo.Categories = cakeCategories
		mu.Unlock()
	}()

	// Получаем начинки торта
	go func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		ids, filIdsErr := u.repo.CakeFillingsIDs(ctx, cakeInfo.ID)
		if filIdsErr != nil {
			trySendError(filIdsErr, errChan, cancel)
			return
		}

		// Получаем информацию каждой категории
		filMu := sync.Mutex{}
		filWg := sync.WaitGroup{}

		var fillings []models.Filling
		for _, fillingID := range ids {
			filWg.Add(1)
			go func() {
				defer filWg.Done()
				if ctx.Err() != nil {
					return
				}

				filling, fillErr := u.repo.FillingByID(ctx, fillingID)
				if fillErr != nil {
					trySendError(fillErr, errChan, cancel)
					return
				}

				filMu.Lock()
				fillings = append(fillings, *filling)
				filMu.Unlock()
			}()
		}

		filWg.Wait()
		if ctx.Err() != nil {
			return
		}

		mu.Lock()
		cakeInfo.Fillings = fillings
		mu.Unlock()
	}()

	// Получаем фотографии торта
	go func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		images, imgErr := u.repo.CakeImages(ctx, cakeInfo.ID)
		if imgErr != nil {
			trySendError(imgErr, errChan, cancel)
			return
		}

		if ctx.Err() != nil {
			return
		}
		mu.Lock()
		cakeInfo.Images = images
		mu.Unlock()
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err = <-errChan; err != nil {
		return nil, err
	}

	return &dto.GetCakeRes{
		Cake: cakeInfo,
	}, nil
}

func (u *CakeUseсase) CreateCake(ctx context.Context, in dto.CreateCakeReq) (*dto.CreateCakeRes, error) {
	// Достаём userID из токена если он не протух
	userID, err := u.tokenator.GetUserIDFromToken(in.AccessToken, false)
	if err != nil {
		return nil, err
	}

	// Добавляем изображение в хранилище
	images := make(map[ms.ImageID][]byte, len(in.Images)+1) // Size это фотографии тортов + превью фотография
	for _, imageData := range in.Images {
		cakeID := ms.ImageID(uuid.New().String())
		images[cakeID] = imageData
	}
	previewImageID := ms.ImageID(uuid.New().String())
	images[previewImageID] = in.PreviewImageData
	res, err := u.imageStore.SaveImages(ctx, u.bucketName, images)
	if err != nil {
		return nil, err
	}

	// Получаем preview
	previewImageURL, ok := res[previewImageID]
	if !ok {
		return nil, errs.ErrPreviewImageNotFound
	}

	// Создаём торт в бд
	cakeID := uuid.New()
	if err = u.repo.CreateCake(ctx, in.ConvertToCreateCakeDBReq(cakeID.String(), previewImageURL, userID, res)); err != nil {
		return nil, err
	}

	return &dto.CreateCakeRes{
		CakeID: cakeID.String(),
	}, nil
}

func (u *CakeUseсase) CreateFilling(ctx context.Context, in dto.CreateFillingReq) (*dto.CreateFillingRes, error) {
	// Достаём userID из токена если он не протух
	_, err := u.tokenator.GetUserIDFromToken(in.AccessToken, false)
	if err != nil {
		return nil, err
	}

	fillingID := uuid.New()
	// Добавляем изображение в хранилище
	imageURL, err := u.imageStore.SaveImage(ctx, u.bucketName, ms.ImageID(fillingID.String()), in.ImageData)
	if err != nil {
		return nil, err
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

	return &dto.CreateFillingRes{
		Filling: filling,
	}, nil
}

func (u *CakeUseсase) CreateCategory(ctx context.Context, in *dto.CreateCategoryReq) (*dto.CreateCategoryRes, error) {
	// Достаём userID из токена если он не протух
	_, err := u.tokenator.GetUserIDFromToken(in.AccessToken, false)
	if err != nil {
		return nil, err
	}

	categoryUUID := uuid.New()
	imageURL, err := u.imageStore.SaveImage(ctx, u.bucketName, ms.ImageID(categoryUUID.String()), in.ImageData)
	if err != nil {
		return nil, err
	}

	newCategory := models.Category{
		ID:       categoryUUID,
		Name:     in.Name,
		ImageURL: imageURL,
	}
	if err = u.repo.CreateCategory(ctx, &newCategory); err != nil {
		return nil, err
	}

	return &dto.CreateCategoryRes{
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

func (u *CakeUseсase) CategoryPreviewCakes(ctx context.Context, categoryID uuid.UUID) ([]*dto.PreviewCake, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Получаем id тортов категории
	cakeIDs, err := u.repo.CategoryCakesIDs(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// Получаем данные тортов по id
	previewCakes := make([]*dto.PreviewCake, len(cakeIDs))
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	errCh := make(chan error, 1)
	for i, cakeID := range cakeIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ctx.Err() != nil {
				return
			}

			previewCake, prevErr := u.repo.PreviewCakeByID(ctx, cakeID)
			if prevErr != nil {
				trySendError(prevErr, errCh, cancel)
				return
			}

			if ctx.Err() != nil {
				return
			}

			mu.Lock()
			previewCakes[i] = previewCake
			mu.Unlock()
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err = <-errCh; err != nil {
		return nil, err
	}

	return previewCakes, nil
}

// trySendError Вспомогательная функция для безопасной отправки ошибки
func trySendError(err error, errCh chan<- error, cancel context.CancelFunc) {
	select {
	case errCh <- err:
		cancel()
	default:
		// Если ошибка уже есть - игнорируем (сохраняем первую)
	}
}
