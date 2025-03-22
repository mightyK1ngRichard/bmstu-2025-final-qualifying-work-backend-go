package repo

import (
	"2025_CakeLand_API/internal/models"
	en "2025_CakeLand_API/internal/pkg/cake/entities"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

const (
	getCakeByID = `
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
	getCakes = `
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
	createFilling = `
		INSERT INTO "filling" (id, name, image_url, content, kg_price, description)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	createCategory = `
		INSERT INTO "category" (id, name, image_url)
		VALUES ($1, $2, $3);
	`
	createCake = `
		INSERT INTO "cake" (id, name, image_url, kg_price, rating, description, mass, is_open_for_sale, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`
	addCakeImages   = `INSERT INTO "cake_images" (id, cake_id, image_url) VALUES ($1, $2, $3);`
	addCateCategory = `
		INSERT INTO "cake_category" (id, category_id, cake_id)
		VALUES ($1, $2, $3);
    `
	addFillingCategory = `
		INSERT INTO "cake_filling" (id, cake_id, filling_id)
		VALUES ($1, $2, $3);
    `
	categories = `SELECT id, name, image_url FROM "category";`
	fillings   = `SELECT id, name, image_url, content, kg_price, description FROM "filling";`
)

type CakeRepository struct {
	db *sql.DB
}

func NewCakeRepository(db *sql.DB) *CakeRepository {
	return &CakeRepository{
		db: db,
	}
}

func (r *CakeRepository) CakeByID(ctx context.Context, in en.GetCakeReq) (*en.GetCakeRes, error) {
	rows, err := r.db.QueryContext(ctx, getCakeByID, in.CakeID)
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
		var dbFilling models.DBFilling
		var category models.Category
		var dbCategory models.DBCategory
		var owner models.User
		var cakeImage models.CakeImage

		// Read data
		if err = rows.Scan(
			&cake.ID, &cake.Name, &cake.ImageURL, &cake.KgPrice, &cake.Rating, &cake.Description, &cake.Mass,
			&cake.IsOpenForSale, &cake.DateCreation, &cake.DiscountKgPrice, &cake.DiscountEndTime,
			&owner.ID, &owner.FIO, &owner.Nickname, &owner.Mail,
			&dbFilling.ID, &dbFilling.Name, &dbFilling.ImageURL, &dbFilling.Content, &dbFilling.KgPrice, &dbFilling.Description,
			&dbCategory.ID, &dbCategory.Name, &dbCategory.ImageURL,
			&cakeImage.ID, &cakeImage.ImageURL,
		); err != nil {
			return nil, models.NewDataBaseError("CakeByID", err)
		}

		if f := dbFilling.ConvertToFilling(); f != nil {
			filling = *f
		}

		if c := dbCategory.ConvertToCategory(); c != nil {
			category = *c
		}

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

	return &en.GetCakeRes{
		Cake: cake,
	}, nil
}

func (r *CakeRepository) CreateCake(ctx context.Context, in en.CreateCakeDBReq) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.NewDataBaseError("CreateCake", err)
	}

	// Создаём торт
	_, err = tx.ExecContext(ctx, createCake,
		in.ID, in.Name, in.PreviewImageURL, in.KgPrice, in.Rating,
		in.Description, in.Mass, in.IsOpenForSale, in.OwnerID,
	)
	if err != nil {
		_ = tx.Rollback()
		return models.NewDataBaseError("[CreateCake]: sqlCommand: createCake", err)
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
		execInGoroutine(addCateCategory, uuid.New(), categoryID, in.ID)
	}

	// Добавляем начинки к торту
	for _, fillingID := range in.FillingIDs {
		execInGoroutine(addFillingCategory, uuid.New(), in.ID, fillingID)
	}

	// Добавляем изображения торта
	for imageID, imageURL := range in.Images {
		execInGoroutine(addCakeImages, imageID, in.ID, imageURL)
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
	if _, err := r.db.ExecContext(ctx, createFilling,
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
	if _, err := r.db.ExecContext(ctx, createCategory, in.ID, in.Name, in.ImageURL); err != nil {
		return models.NewDataBaseError("CreateCategory", err)
	}

	return nil
}

func (r *CakeRepository) Categories(ctx context.Context) (*[]models.Category, error) {
	rows, err := r.db.QueryContext(ctx, categories)
	defer rows.Close()
	if err != nil {
		return nil, models.NewDataBaseError("Categories", err)
	}

	// Чтение результатов
	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err = rows.Scan(&category.ID, &category.Name, &category.ImageURL); err != nil {
			return nil, models.NewDataBaseError("Categories", err)
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
	rows, err := r.db.QueryContext(ctx, fillings)
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
	rows, err := r.db.QueryContext(ctx, getCakes)
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
		var dbFilling models.DBFilling
		var category models.Category
		var dbCategory models.DBCategory
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

		if f := dbFilling.ConvertToFilling(); f != nil {
			filling = *f
		}

		if c := dbCategory.ConvertToCategory(); c != nil {
			category = *c
		}

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

// Функция для преобразования map[uuid.UUID]Cake в []Cake
func mapToSlice(cakes map[uuid.UUID]models.Cake) []models.Cake {
	result := make([]models.Cake, 0, len(cakes))
	for _, cake := range cakes {
		result = append(result, cake)
	}
	return result
}
