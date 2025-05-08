package usecase

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	dto2 "2025_CakeLand_API/internal/pkg/cake/dto"
	"2025_CakeLand_API/internal/pkg/minio"
	"2025_CakeLand_API/internal/pkg/profile"
	gen "2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/profile/dto"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	"context"
	"github.com/google/uuid"
	"strings"
	"sync"
)

type ProfileUseсase struct {
	tokenator     *jwt.Tokenator
	repo          profile.IProfileRepository
	imageProvider *minio.MinioProvider
}

func NewProfileUsecase(
	tokenator *jwt.Tokenator,
	repo profile.IProfileRepository,
	imageProvider *minio.MinioProvider,
) *ProfileUseсase {
	return &ProfileUseсase{
		tokenator:     tokenator,
		repo:          repo,
		imageProvider: imageProvider,
	}
}

func (u *ProfileUseсase) UpdateUserAddresses(ctx context.Context, accessToken string, req *gen.UpdateUserAddressesReq) (models.Address, error) {
	// Достаём UserID
	userID, err := u.getUserUUID(accessToken)
	if err != nil {
		return models.Address{}, err
	}

	// Обновляем адрес в репозитории
	return u.repo.UpdateUserAddresses(ctx, userID, req)
}

func (u *ProfileUseсase) GetUserAddresses(ctx context.Context, accessToken string) ([]models.Address, error) {
	// Достаём UserID
	userID, err := u.getUserUUID(accessToken)
	if err != nil {
		return nil, err
	}

	// Получаем адреса из БД
	return u.repo.GetUserAddresses(ctx, userID)
}

func (u *ProfileUseсase) CreateAddress(ctx context.Context, accessToken string, address *models.Address) (*models.Address, error) {
	// Достаём UserID
	userID, err := u.getUserUUID(accessToken)
	if err != nil {
		return nil, err
	}

	// Создаём адрес в БД
	address.ID = uuid.New()
	address.UserID = userID

	if err = u.repo.CreateAddress(ctx, address); err != nil {
		return nil, err
	}

	// Ответ
	return address, err
}

func (u *ProfileUseсase) UserInfo(ctx context.Context, accessToken string) (*dto.UserInfo, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Достаём UserID
	userID, err := u.getUserUUID(accessToken)
	if err != nil {
		return nil, err
	}

	//var userInfo dto.UserInfo
	var (
		dbCakes []dto2.PreviewCakeDB
		user    *dto.Profile
	)
	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	wg.Add(2)

	// Получаем данные пользователя
	go func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		profileDB, err := u.repo.UserInfo(ctx, userID)
		if err != nil {
			trySendError(err, errChan, cancel)
			return
		}

		if ctx.Err() != nil {
			return
		}

		user = profileDB
	}()

	// Получаем тортики пользователя
	go func() {
		defer wg.Done()
		if ctx.Err() != nil {
			return
		}

		dbPreviewCakes, err := u.repo.CakesByUserID(ctx, userID)
		if err != nil {
			trySendError(err, errChan, cancel)
			return
		}

		if ctx.Err() != nil {
			return
		}

		dbCakes = dbPreviewCakes
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err = <-errChan; err != nil {
		return nil, err
	}

	cakes := make([]dto2.PreviewCake, len(dbCakes))
	owner := user.ConvertToOwner()
	for i, cake := range dbCakes {
		cakes[i] = cake.ConvertToPreviewCake(owner)
	}
	userInfo := dto.UserInfo{
		User:  *user,
		Cakes: cakes,
	}
	return &userInfo, nil
}

func (u *ProfileUseсase) UserInfoByID(ctx context.Context, userID uuid.UUID) (*models.UserInfo, error) {
	profileInfo, err := u.repo.UserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	userInfo := models.UserInfo{
		ID:             profileInfo.ID.String(),
		FIO:            profileInfo.FIO,
		Nickname:       profileInfo.Nickname,
		ImageURL:       profileInfo.ImageURL,
		HeaderImageURL: profileInfo.HeaderImageURL,
		Mail:           profileInfo.Mail,
		Phone:          profileInfo.Phone,
	}
	return &userInfo, nil
}

func (u *ProfileUseсase) UpdateUserImage(ctx context.Context, accessToken string, in *gen.UpdateUserImageReq) (string, error) {
	// Получаем ID пользователя
	userUUID, err := u.getUserUUID(accessToken)
	if err != nil {
		return "", err
	}

	// Созхраняем фото в минио
	imageID := uuid.New()
	imageURL, err := u.imageProvider.SaveImage(ctx, minio.ImageID(imageID.String()), in.ImageData)
	if err != nil {
		return "", err
	}

	// Обновляем запись в БД
	switch in.ImageKind {
	case gen.UpdateUserImageReq_AVATAR:
		return imageURL, u.repo.UpdateUserAvatar(ctx, userUUID, imageURL)
	case gen.UpdateUserImageReq_HEADER:
		return imageURL, u.repo.UpdateUserHeaderImage(ctx, userUUID, imageURL)
	default:
		return "", errs.ErrUnknownImageKind
	}
}

func (u *ProfileUseсase) UpdateUserData(ctx context.Context, accessToken string, in *gen.UpdateUserDataReq) error {
	// Получаем UseriD
	userUUID, err := u.getUserUUID(accessToken)
	if err != nil {
		return err
	}

	// Убираем лишние пробелы в имени и ФИО
	in.UpdatedUserName = strings.TrimSpace(in.UpdatedUserName)
	if in.UpdatedFIO != nil {
		*in.UpdatedFIO = strings.TrimSpace(*in.UpdatedFIO)
	}

	// Сохраняем данные в БД
	return u.repo.UpdateUserData(ctx, userUUID, in)
}

func trySendError(err error, errCh chan<- error, cancel context.CancelFunc) {
	select {
	case errCh <- err:
		cancel()
	default:
		// Если ошибка уже есть - игнорируем (сохраняем первую)
	}
}

func (u *ProfileUseсase) getUserUUID(accessToken string) (uuid.UUID, error) {
	// Достаём UserID
	userIDStr, err := u.tokenator.GetUserIDFromToken(accessToken, false)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
