package repo

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/cake/dto"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"sync"
)

const (
	queryGetCakeByID = `
        SELECT c.id, c.name, c.image_url, c.kg_price, c.rating, c.description, c.mass, c.is_open_for_sale, 
               c.date_creation, c.discount_kg_price, c.discount_end_time,
		   u.id AS owner_id, u.fio, u.nickname, u.mail,
		   f.id AS filling_id, f.name AS filling_name, f.image_url AS filling_image,
		   f.content AS filling_content, f.kg_price AS filling_price_per_kg, f.description AS filling_description,
		   cat.id AS category_id, cat.name AS category_name, cat.image_url AS category_image,
		   ci.id, ci.image_url
        FROM "cake" c
			 LEFT JOIN "user" u ON c.owner_id = u.id
			 LEFT JOIN "cake_filling" cf ON c.id = cf.cake_id
			 LEFT JOIN "filling" f ON cf.filling_id = f.id
			 LEFT JOIN "cake_category" cc ON c.id = cc.cake_id
			 LEFT JOIN "category" cat ON cc.category_id = cat.id
			 LEFT JOIN "cake_images" ci ON ci.cake_id = c.id
        WHERE c.id = $1;
    `
	queryGetCakes = `
        SELECT c.id, c.name, c.image_url, c.kg_price, c.rating, c.description, c.mass, c.is_open_for_sale,
               c.date_creation, c.discount_kg_price, c.discount_end_time,
               u.id AS owner_id, u.fio, u.nickname, u.mail,
               f.id AS filling_id, f.name AS filling_name, f.image_url AS filling_image,
               f.content AS filling_content, f.kg_price AS filling_price_per_kg, f.description AS filling_description,
               cat.id AS category_id, cat.name AS category_name, cat.image_url AS category_image
        FROM "cake" c
                 LEFT JOIN "user" u ON c.owner_id = u.id
                 LEFT JOIN "cake_filling" cf ON c.id = cf.cake_id
                 LEFT JOIN "filling" f ON cf.filling_id = f.id
                 LEFT JOIN "cake_category" cc ON c.id = cc.cake_id
                 LEFT JOIN "category" cat ON cc.category_id = cat.id
    `
	queryCreateFilling = `
		INSERT INTO "filling" (id, name, image_url, content, kg_price, description)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	queryCreateCategory = `
		INSERT INTO "category" (id, name, image_url)
		VALUES ($1, $2, $3);
	`
	queryCreateCake = `
		INSERT INTO "cake" (id, name, image_url, kg_price, rating, description, mass, is_open_for_sale, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
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
		SELECT id,
			   name,
			   image_url,
			   kg_price,
			   rating,
			   description,
			   mass,
			   discount_kg_price,
			   discount_end_time,
			   date_creation,
			   is_open_for_sale,
			   owner_id
		FROM cake
		WHERE id = $1;
	`
)

type CakeRepository struct {
	db *sql.DB
}

func NewCakeRepository(db *sql.DB) *CakeRepository {
	return &CakeRepository{
		db: db,
	}
}

func (r *CakeRepository) CakeByID(ctx context.Context, in dto.GetCakeReq) (*dto.GetCakeRes, error) {
	rows, err := r.db.QueryContext(ctx, queryGetCakeByID, in.CakeID)
	if err != nil {
		return nil, models.NewDataBaseError("CakeByID", err)
	}

	defer rows.Close()
	var cake models.Cake
	cake.Fillings = []models.Filling{}
	cake.Categories = []models.Category{}
	cake.Images = []models.CakeImage{}

	// Map for unique fillings and categories to avoid duplicates
	fillingMap := make(map[uuid.UUID]bool)
	categoryMap := make(map[uuid.UUID]bool)
	cakeImagesMap := make(map[uuid.UUID]bool)

	for rows.Next() {
		var filling models.Filling
		var dbFilling dto.DBFilling
		var category models.Category
		var dbCategory dto.DBCategory
		var owner models.User
		var cakeImage models.CakeImage

		// Read data
		if err = rows.Scan(
			&cake.ID, &cake.Name, &cake.ImageURL, &cake.KgPrice, &cake.Rating, &cake.Description, &cake.Mass,
			&cake.IsOpenForSale, &cake.DateCreation, &cake.DiscountKgPrice, &cake.DiscountEndTime,
			&owner.ID, &owner.FIO, &owner.Nickname, &owner.Mail,
			&dbFilling.ID, &dbFilling.Name, &dbFilling.ImageURL, &dbFilling.Content, &dbFilling.KgPrice, &dbFilling.Description,
			&dbCategory.ID, &dbCategory.Name, &dbCategory.ImageURL, &dbCategory.CategoryGenders,
			&cakeImage.ID, &cakeImage.ImageURL,
		); err != nil {
			return nil, models.NewDataBaseError("CakeByID", err)
		}

		filling = dbFilling.ConvertToFilling()
		category = dbCategory.ConvertToCategory()

		// Set owner (only once, since it's the same for all rows)
		if cake.Owner.ID == uuid.Nil {
			cake.Owner = owner
		}

		// Add unique fillings
		if filling.ID != uuid.Nil && !fillingMap[filling.ID] {
			fillingMap[filling.ID] = true
			cake.Fillings = append(cake.Fillings, filling)
		}

		// Add unique categories
		if category.ID != uuid.Nil && !categoryMap[category.ID] {
			categoryMap[category.ID] = true
			cake.Categories = append(cake.Categories, category)
		}

		// Add unique cake image
		if cakeImage.ID != uuid.Nil && !cakeImagesMap[cakeImage.ID] {
			cakeImagesMap[cakeImage.ID] = true
			cake.Images = append(cake.Images, cakeImage)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, models.NewDataBaseError("CakeByID", err)
	}

	return &dto.GetCakeRes{
		Cake: cake,
	}, nil
}

func (r *CakeRepository) CreateCake(ctx context.Context, in dto.CreateCakeDBReq) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.NewDataBaseError("CreateCake", err)
	}

	// Создаём торт
	_, err = tx.ExecContext(ctx, queryCreateCake,
		in.ID, in.Name, in.PreviewImageURL, in.KgPrice, in.Rating,
		in.Description, in.Mass, in.IsOpenForSale, in.OwnerID,
	)
	if err != nil {
		_ = tx.Rollback()
		return models.NewDataBaseError("[CreateCake]: sqlCommand: queryCreateCake", err)
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
				case errChan <- models.NewDataBaseError(fmt.Sprintf("[CreateCake]: %s", query), dbErr):
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
		return dbErr
	case <-done:
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return models.NewDataBaseError("CreateCake 4", err)
	}

	return nil
}

func (r *CakeRepository) CreateFilling(ctx context.Context, in models.Filling) error {
	if _, err := r.db.ExecContext(ctx, queryCreateFilling,
		in.ID,
		in.Name,
		in.ImageURL,
		in.Content,
		in.KgPrice,
		in.Description,
	); err != nil {
		return models.NewDataBaseError("CreateFilling", err)
	}

	return nil
}

func (r *CakeRepository) CreateCategory(ctx context.Context, in *models.Category) error {
	// TODO: Сделать добавление тегов
	if _, err := r.db.ExecContext(ctx, queryCreateCategory, in.ID, in.Name, in.ImageURL); err != nil {
		return models.NewDataBaseError("CreateCategory", err)
	}

	return nil
}

func (r *CakeRepository) Categories(ctx context.Context) (*[]models.Category, error) {
	rows, err := r.db.QueryContext(ctx, queryCategories)
	defer rows.Close()
	if err != nil {
		return nil, models.NewDataBaseError("Categories", err)
	}

	// Чтение результатов
	var categories []models.Category
	for rows.Next() {
		var (
			category        models.Category
			categoryGenders []string
		)

		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.ImageURL,
			pq.Array(&categoryGenders),
		); err != nil {
			return nil, models.NewDataBaseError("Categories", err)
		}

		// Преобразуем []string в []CategoryGender
		category.CategoryGenders = make([]models.CategoryGender, 0, len(categoryGenders))
		for _, g := range categoryGenders {
			category.CategoryGenders = append(category.CategoryGenders, models.CategoryGender(g))
		}

		categories = append(categories, category)
	}

	// Проверка на ошибки после завершения итерации
	if err = rows.Err(); err != nil {
		return nil, models.NewDataBaseError("Categories", err)
	}

	return &categories, nil
}

func (r *CakeRepository) Fillings(ctx context.Context) (*[]models.Filling, error) {
	rows, err := r.db.QueryContext(ctx, queryFillings)
	defer rows.Close()
	if err != nil {
		return nil, models.NewDataBaseError("Fillings", err)
	}

	// Чтение результатов
	var fillings []models.Filling
	for rows.Next() {
		var filling models.Filling
		if err = rows.Scan(
			&filling.ID, &filling.Name, &filling.ImageURL,
			&filling.Content, &filling.KgPrice, &filling.Description,
		); err != nil {
			return nil, models.NewDataBaseError("Fillings", err)
		}
		fillings = append(fillings, filling)
	}

	// Проверка на ошибки после завершения итерации
	if err = rows.Err(); err != nil {
		return nil, models.NewDataBaseError("Fillings", err)
	}

	return &fillings, nil
}

func (r *CakeRepository) Cakes(ctx context.Context) (*[]models.Cake, error) {
	rows, err := r.db.QueryContext(ctx, queryGetCakes)
	if err != nil {
		return nil, models.NewDataBaseError("Cakes", err)
	}
	defer rows.Close()

	// Map для уникальных
	cakes := make(map[uuid.UUID]models.Cake)
	fillingMap := make(map[string]bool)
	categoryMap := make(map[string]bool)

	// Обработка строк в запросе
	for rows.Next() {
		var cake models.Cake
		var filling models.Filling
		var dbFilling dto.DBFilling
		var category models.Category
		var dbCategory dto.DBCategory
		var owner models.User

		// Чтение данных
		if err = rows.Scan(
			&cake.ID, &cake.Name, &cake.ImageURL, &cake.KgPrice, &cake.Rating, &cake.Description, &cake.Mass,
			&cake.IsOpenForSale, &cake.DateCreation, &cake.DiscountKgPrice, &cake.DiscountEndTime,
			&owner.ID, &owner.FIO, &owner.Nickname, &owner.Mail,
			&dbFilling.ID, &dbFilling.Name, &dbFilling.ImageURL, &dbFilling.Content, &dbFilling.KgPrice, &dbFilling.Description,
			&dbCategory.ID, &dbCategory.Name, &dbCategory.ImageURL,
		); err != nil {
			return nil, models.NewDataBaseError("Cakes", err)
		}

		filling = dbFilling.ConvertToFilling()
		category = dbCategory.ConvertToCategory()

		// Достаём торт если он есть или инициализируем отсканированный
		savedCake, ok := cakes[cake.ID]
		if !ok {
			savedCake = cake
		}

		// Устанавливаем владельца только один раз, так как он одинаковый для всех строк
		if savedCake.Owner.ID == uuid.Nil {
			savedCake.Owner = owner
		}

		// Добавляем уникальные начинки
		// Создаём ключ из ID торта и ID начинки для уникальности начинок для каждого торта
		key := savedCake.ID.String() + filling.ID.String()
		if filling.ID != uuid.Nil && !fillingMap[key] {
			fillingMap[key] = true
			savedCake.Fillings = append(savedCake.Fillings, filling)
		}

		// Добавляем уникальные категории
		key = savedCake.ID.String() + category.ID.String()
		if category.ID != uuid.Nil && !categoryMap[key] {
			categoryMap[key] = true
			savedCake.Categories = append(savedCake.Categories, category)
		}

		// Добавляем торт в общий список
		cakes[savedCake.ID] = savedCake
	}

	if err = rows.Err(); err != nil {
		return nil, models.NewDataBaseError("Cakes", err)
	}

	cakeSlice := mapToSlice(cakes)
	return &cakeSlice, nil
}

func (r *CakeRepository) CategoryIDsByGenderName(ctx context.Context, genderTag models.CategoryGender) ([]dto.DBCategory, error) {
	rows, err := r.db.QueryContext(ctx, queryCakesByGenderTag, genderTag)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var categories []dto.DBCategory
	for rows.Next() {
		var category dto.DBCategory
		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.ImageURL,
			pq.Array(&category.CategoryGenders),
		); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CakeRepository) CategoryCakesIDs(ctx context.Context, categoryID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.QueryContext(ctx, queryCategoryCakesIDs, categoryID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *CakeRepository) PreviewCakeByID(ctx context.Context, cakeID uuid.UUID) (*dto.PreviewCake, error) {
	var dbPreviewCake dto.PreviewCake
	if err := r.db.QueryRowContext(ctx, queryPreviewCakeByID, cakeID).Scan(
		&dbPreviewCake.ID,
		&dbPreviewCake.Name,
		&dbPreviewCake.PreviewImageURL,
		&dbPreviewCake.KgPrice,
		&dbPreviewCake.Rating,
		&dbPreviewCake.Description,
		&dbPreviewCake.Mass,
		&dbPreviewCake.DiscountKgPrice,
		&dbPreviewCake.DiscountEndTime,
		&dbPreviewCake.DateCreation,
		&dbPreviewCake.IsOpenForSale,
		&dbPreviewCake.OwnerID,
	); err != nil {
		return nil, err
	}

	return &dbPreviewCake, nil
}

// Функция для преобразования map[uuid.UUID]Cake в []Cake
func mapToSlice(cakes map[uuid.UUID]models.Cake) []models.Cake {
	result := make([]models.Cake, 0, len(cakes))
	for _, cake := range cakes {
		result = append(result, cake)
	}
	return result
}
