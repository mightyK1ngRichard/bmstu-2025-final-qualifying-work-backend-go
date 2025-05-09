package repo

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/cake/dto"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"sync"
)

const (
	queryGetCakeCategoriesIDs = `SELECT category_id FROM cake_category WHERE cake_id = $1`
	queryGetCakeFillingsIDs   = `SELECT filling_id FROM cake_filling WHERE cake_id = $1`
	queryGetFillingByID       = `SELECT id, name, image_url, content, kg_price, description FROM filling WHERE id = $1`
	queryGetCategoryByID      = `SELECT id, name, image_url, gender_tags FROM category WHERE id = $1`
	queryGetCakeImages        = `SELECT id, image_url FROM cake_images WHERE cake_id = $1`
	queryGetCakeByID          = `
		SELECT c.id, c.name, c.image_url, c.kg_price, c.reviews_count, c.stars_sum,
			   c.description, c.mass, c.is_open_for_sale, c.date_creation, c.discount_kg_price, c.discount_end_time,
			   u.id AS owner_id, u.fio, u.address, u.nickname, u.image_url, u.mail, u.phone, u.header_image_url
		FROM "cake" c
				 LEFT JOIN "user" u ON c.owner_id = u.id
		WHERE c.id = $1
	`
	queryCreateFilling = `
		INSERT INTO "filling" (id, name, image_url, content, kg_price, description)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	queryCreateCategory = `INSERT INTO "category" (id, name, image_url) VALUES ($1, $2, $3);`
	queryCreateCake     = `
		INSERT INTO "cake" (id, name, image_url, kg_price, discount_kg_price, discount_end_time, description, mass, is_open_for_sale, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`
	queryAddCakeImages   = `INSERT INTO "cake_images" (id, cake_id, image_url) VALUES ($1, $2, $3);`
	queryAddCateCategory = `
		INSERT INTO "cake_category" (id, category_id, cake_id)
		VALUES ($1, $2, $3);
    `
	queryAddFillingCategory = `
		INSERT INTO "cake_filling" (id, cake_id, filling_id)
		VALUES ($1, $2, $3);
    `
	queryCategories       = `SELECT id, name, image_url, gender_tags FROM "category";`
	queryFillings         = `SELECT id, name, image_url, content, kg_price, description FROM "filling";`
	queryCakesByGenderTag = `SELECT id, name, image_url, gender_tags FROM category WHERE $1 = ANY(gender_tags);`
	queryCategoryCakesIDs = `SELECT cake_id FROM cake_category WHERE category_id = $1;`
	queryPreviewCakeByID  = `
		SELECT c.id,
			   c.name,
			   c.image_url,
			   c.kg_price,
			   c.reviews_count,
			   c.stars_sum,
			   c.description,
			   c.mass,
			   c.discount_kg_price,
			   c.discount_end_time,
			   c.date_creation,
			   c.is_open_for_sale,
			   u.id,
			   u.fio,
			   u.address,
			   u.nickname,
			   u.image_url,
			   u.mail,
			   u.phone,
			   u.header_image_url
		FROM cake c
				 LEFT JOIN "user" u ON u.id = c.owner_id
		WHERE c.id = $1
	`
	queryGetColors    = `SELECT DISTINCT hex_color FROM cake_color`
	queryAddCakeColor = `INSERT INTO cake_color (id, cake_id, hex_color) VALUES ($1, $2, $3)`
	queryGetAllCakes  = `
		SELECT id,
			   name,
			   image_url,
			   kg_price,
			   reviews_count,
			   stars_sum,
			   description,
			   mass,
			   discount_kg_price,
			   discount_end_time,
			   date_creation,
			   is_open_for_sale,
			   owner_id
		FROM cake
	`
	queryGetUser = `
		SELECT id,
			   fio,
			   address,
			   nickname,
			   image_url,
			   header_image_url,
			   mail,
			   phone
		FROM "user"
		WHERE id = $1
	`
	queryGetCakeColors = `SELECT id, cake_id, hex_color FROM cake_color WHERE cake_id = $1`
)

type CakeRepository struct {
	db *sql.DB
}

func NewCakeRepository(db *sql.DB) *CakeRepository {
	return &CakeRepository{
		db: db,
	}
}

func (r *CakeRepository) GetCakeColorsByCakeID(ctx context.Context, cakeID uuid.UUID) ([]models.CakeColor, error) {
	const methodName = "[CakeRepository.GetCakeColorsByCakeID]"

	rows, err := r.db.QueryContext(ctx, queryGetCakeColors, cakeID)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}
	defer rows.Close()

	var colors []models.CakeColor
	for rows.Next() {
		var color models.CakeColor
		if err = rows.Scan(&color.ID, &color.CakeID, &color.HexString); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}

		colors = append(colors, color)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return colors, nil
}

func (r *CakeRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (dto.Owner, error) {
	const methodName = "[UserRepository.GetUserByID]"

	row := r.db.QueryRowContext(ctx, queryGetUser, userID)

	var user dto.Owner
	if err := row.Scan(
		&user.ID,
		&user.FIO,
		&user.Address,
		&user.Nickname,
		&user.ImageURL,
		&user.HeaderImageURL,
		&user.Mail,
		&user.Phone,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errs.WrapDBError(methodName, err)
		}
		return user, errs.WrapDBError(methodName, err)
	}

	return user, nil
}

func (r *CakeRepository) GetCakesPreview(ctx context.Context) ([]dto.PreviewCake, error) {
	const methodName = "[CakeRepository.GetCakes]"

	rows, err := r.db.QueryContext(ctx, queryGetAllCakes)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}
	defer rows.Close()

	var cakes []dto.PreviewCake
	for rows.Next() {
		var cake dto.PreviewCake
		var discountKgPrice sql.NullFloat64
		var discountEndTime sql.NullTime
		var ownerID uuid.UUID

		if err = rows.Scan(
			&cake.ID,
			&cake.Name,
			&cake.PreviewImageURL,
			&cake.KgPrice,
			&cake.ReviewsCount,
			&cake.StarsSum,
			&cake.Description,
			&cake.Mass,
			&discountKgPrice,
			&discountEndTime,
			&cake.DateCreation,
			&cake.IsOpenForSale,
			&ownerID,
		); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}

		cake.DiscountKgPrice = null.FloatFromPtr(nil)
		if discountKgPrice.Valid {
			cake.DiscountKgPrice = null.FloatFrom(discountKgPrice.Float64)
		}

		cake.DiscountEndTime = null.TimeFromPtr(nil)
		if discountEndTime.Valid {
			cake.DiscountEndTime = null.TimeFrom(discountEndTime.Time)
		}

		cake.Owner.ID = ownerID
		cakes = append(cakes, cake)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return cakes, nil
}

func (r *CakeRepository) AddCakeColor(ctx context.Context, in models.CakeColor) error {
	const methodName = "[CakeRepository.AddCakeColor]"

	if _, err := r.db.ExecContext(ctx, queryAddCakeColor, in.ID, in.CakeID, in.HexString); err != nil {
		return errs.WrapDBError(methodName, err)
	}

	return nil
}

func (r *CakeRepository) CakeByID(ctx context.Context, in dto.GetCakeReq) (*dto.GetCakeRes, error) {
	const methodName = "[CakeRepository.CakeByID]"

	var cake models.Cake
	if err := r.db.QueryRowContext(ctx, queryGetCakeByID, in.CakeID).Scan(
		&cake.ID, &cake.Name, &cake.PreviewImageURL, &cake.KgPrice, &cake.ReviewsCount, &cake.StarsSum, &cake.Description,
		&cake.Mass, &cake.IsOpenForSale, &cake.DateCreation, &cake.DiscountKgPrice, &cake.DiscountEndTime,
		&cake.Owner.ID, &cake.Owner.FIO, &cake.Owner.Address,
		&cake.Owner.Nickname, &cake.Owner.ImageURL, &cake.Owner.Mail, &cake.Owner.Phone,
		&cake.Owner.HeaderImageURL,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.WrapDBError(methodName, err)
	}

	return &dto.GetCakeRes{
		Cake: cake,
	}, nil
}

func (r *CakeRepository) GetColors(ctx context.Context) ([]string, error) {
	const methodName = "[CakeRepository.GetColors]"

	rows, err := r.db.QueryContext(ctx, queryGetColors)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}
	defer rows.Close()

	hexStrings := make([]string, 0, 22)
	for rows.Next() {
		var hexString string
		if err = rows.Scan(&hexString); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}

		hexStrings = append(hexStrings, hexString)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return hexStrings, nil
}

func (r *CakeRepository) CakeCategoriesIDs(ctx context.Context, cakeID uuid.UUID) ([]uuid.UUID, error) {
	const methodName = "[Repo.CakeCategoriesIDs]"

	rows, err := r.db.QueryContext(ctx, queryGetCakeCategoriesIDs, cakeID)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}

		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return ids, nil
}

func (r *CakeRepository) CakeFillingsIDs(ctx context.Context, cakeID uuid.UUID) ([]uuid.UUID, error) {
	const methodName = "[Repo.CakeFillingsIDs]"

	rows, err := r.db.QueryContext(ctx, queryGetCakeFillingsIDs, cakeID)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}

		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return ids, nil
}

func (r *CakeRepository) FillingByID(ctx context.Context, fillingID uuid.UUID) (*models.Filling, error) {
	const methodName = "[Repo.FillingByID]"

	var filling models.Filling
	if err := r.db.QueryRowContext(ctx, queryGetFillingByID, fillingID).Scan(
		&filling.ID,
		&filling.Name,
		&filling.ImageURL,
		&filling.Content,
		&filling.KgPrice,
		&filling.Description,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}

		return nil, errs.WrapDBError(methodName, err)
	}

	return &filling, nil
}

func (r *CakeRepository) CategoryByID(ctx context.Context, categoryID uuid.UUID) (*models.Category, error) {
	const methodName = "[Repo.CategoryByID]"

	var category models.Category
	var genderTags pq.StringArray

	if err := r.db.QueryRowContext(ctx, queryGetCategoryByID, categoryID).Scan(
		&category.ID,
		&category.Name,
		&category.ImageURL,
		&genderTags,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.WrapDBError(methodName, err)
	}

	category.CategoryGenders = models.ParseGenderTags(genderTags)
	return &category, nil
}

func (r *CakeRepository) CakeImages(ctx context.Context, cakeID uuid.UUID) ([]models.CakeImage, error) {
	const methodName = "[Repo.CakeImages]"

	rows, err := r.db.QueryContext(ctx, queryGetCakeImages, cakeID)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	defer rows.Close()
	var images []models.CakeImage
	for rows.Next() {
		var image models.CakeImage
		if err = rows.Scan(&image.ID, &image.ImageURL); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}

		images = append(images, image)
	}

	return images, nil
}

func (r *CakeRepository) CreateCake(ctx context.Context, in dto.CreateCakeDBReq) error {
	const methodName = "[Repo.CreateCake]"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.WrapDBError(methodName, err)
	}

	// Создаём торт
	if _, err = tx.ExecContext(ctx, queryCreateCake,
		in.ID, in.Name, in.PreviewImageURL, in.KgPrice, in.DiscountedKgPrice, in.DiscountedPriceEndDate,
		in.Description, in.Mass, in.IsOpenForSale, in.OwnerID,
	); err != nil {
		_ = tx.Rollback()
		return errs.WrapDBError(methodName, err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	cancelableCtx, cancel := context.WithCancel(ctx)

	// Обёртка для запуска SQL-запросов в горутине
	execInGoroutine := func(query string, args ...interface{}) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, dbErr := tx.ExecContext(cancelableCtx, query, args...)
			if dbErr != nil {
				select {
				case errChan <- dbErr:
					cancel()
				default:
				}
			}
		}()
	}

	// Добавляем категории к торту
	for _, categoryID := range in.CategoryIDs {
		execInGoroutine(queryAddCateCategory, uuid.New(), categoryID, in.ID)
	}

	// Добавляем начинки к торту
	for _, fillingID := range in.FillingIDs {
		execInGoroutine(queryAddFillingCategory, uuid.New(), in.ID, fillingID)
	}

	// Добавляем изображения торта
	for imageID, imageURL := range in.Images {
		execInGoroutine(queryAddCakeImages, imageID, in.ID, imageURL)
	}

	// Ожидаем ошибку или завершение всех горутин
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
		close(errChan)
	}()

	select {
	case dbErr := <-errChan:
		_ = tx.Rollback()
		return errs.WrapDBError(methodName, dbErr)
	case <-done:
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return errs.WrapDBError(methodName, err)
	}

	return nil
}

func (r *CakeRepository) CreateFilling(ctx context.Context, in models.Filling) error {
	const methodName = "[Repo.CreateFilling]"

	if _, err := r.db.ExecContext(ctx, queryCreateFilling,
		in.ID,
		in.Name,
		in.ImageURL,
		in.Content,
		in.KgPrice,
		in.Description,
	); err != nil {
		return errs.WrapDBError(methodName, err)
	}

	return nil
}

func (r *CakeRepository) CreateCategory(ctx context.Context, in *models.Category) error {
	const methodName = "[Repo.CreateCategory]"

	// TODO: Сделать добавление тегов
	if _, err := r.db.ExecContext(ctx, queryCreateCategory, in.ID, in.Name, in.ImageURL); err != nil {
		return errs.WrapDBError(methodName, err)
	}

	return nil
}

func (r *CakeRepository) Categories(ctx context.Context) (*[]models.Category, error) {
	const methodName = "[Repo.Categories]"

	rows, err := r.db.QueryContext(ctx, queryCategories)
	defer rows.Close()
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	// Чтение результатов
	var categories []models.Category
	for rows.Next() {
		var (
			category        models.Category
			categoryGenders pq.StringArray
		)

		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.ImageURL,
			&categoryGenders,
		); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}

		category.CategoryGenders = models.ParseGenderTags(categoryGenders)
		categories = append(categories, category)
	}

	// Проверка на ошибки после завершения итерации
	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return &categories, nil
}

func (r *CakeRepository) Fillings(ctx context.Context) (*[]models.Filling, error) {
	const methodName = "[Repo.Fillings]"

	rows, err := r.db.QueryContext(ctx, queryFillings)
	defer rows.Close()
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	// Чтение результатов
	var fillings []models.Filling
	for rows.Next() {
		var filling models.Filling
		if err = rows.Scan(
			&filling.ID, &filling.Name, &filling.ImageURL,
			&filling.Content, &filling.KgPrice, &filling.Description,
		); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}
		fillings = append(fillings, filling)
	}

	// Проверка на ошибки после завершения итерации
	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return &fillings, nil
}

func (r *CakeRepository) CategoryIDsByGenderName(ctx context.Context, genderTag models.CategoryGender) ([]dto.DBCategory, error) {
	const methodName = "[Repo.CategoryIDsByGender]"

	rows, err := r.db.QueryContext(ctx, queryCakesByGenderTag, genderTag)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	defer rows.Close()
	var categories []dto.DBCategory
	for rows.Next() {
		var (
			category   dto.DBCategory
			genderTags pq.StringArray
		)

		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.ImageURL,
			&genderTags,
		); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}

		category.CategoryGenders = models.ParseGenderTags(genderTags)
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return categories, nil
}

func (r *CakeRepository) CategoryCakesIDs(ctx context.Context, categoryID uuid.UUID) ([]uuid.UUID, error) {
	const methodName = "[Repo.CategoryCakesID]"

	rows, err := r.db.QueryContext(ctx, queryCategoryCakesIDs, categoryID)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}

		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return ids, nil
}

func (r *CakeRepository) PreviewCakeByID(ctx context.Context, cakeID uuid.UUID) (*dto.PreviewCake, error) {
	const methodName = "[Repo.PreviewCakeByID]"

	var previewCake dto.PreviewCake
	if err := r.db.QueryRowContext(ctx, queryPreviewCakeByID, cakeID).Scan(
		&previewCake.ID,
		&previewCake.Name,
		&previewCake.PreviewImageURL,
		&previewCake.KgPrice,
		&previewCake.ReviewsCount,
		&previewCake.StarsSum,
		&previewCake.Description,
		&previewCake.Mass,
		&previewCake.DiscountKgPrice,
		&previewCake.DiscountEndTime,
		&previewCake.DateCreation,
		&previewCake.IsOpenForSale,
		&previewCake.Owner.ID,
		&previewCake.Owner.FIO,
		&previewCake.Owner.Address,
		&previewCake.Owner.Nickname,
		&previewCake.Owner.ImageURL,
		&previewCake.Owner.Mail,
		&previewCake.Owner.Phone,
		&previewCake.Owner.HeaderImageURL,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.WrapDBError(methodName, err)
	}

	return &previewCake, nil
}

// Функция для преобразования map[uuid.UUID]Cake в []Cake
func mapToSlice(cakes map[uuid.UUID]models.Cake) []models.Cake {
	result := make([]models.Cake, 0, len(cakes))
	for _, cake := range cakes {
		result = append(result, cake)
	}
	return result
}
