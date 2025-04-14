package repo

import (
	"2025_CakeLand_API/internal/models/errs"
	cakeDto "2025_CakeLand_API/internal/pkg/cake/dto"
	"2025_CakeLand_API/internal/pkg/profile/dto"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

const (
	querySelectProfileByID   = `SELECT id, fio, address, nickname, header_image_url, image_url, mail, phone, card_number FROM "user" WHERE id = $1 LIMIT 1`
	querySelectCakesByUserID = `
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
		WHERE owner_id = $1;
    `
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) UserInfo(ctx context.Context, userID uuid.UUID) (*dto.Profile, error) {
	methodName := "[Repo.UserInfo]"

	var user dto.Profile
	if err := r.db.QueryRowContext(ctx, querySelectProfileByID, userID).Scan(
		&user.ID,
		&user.FIO,
		&user.Address,
		&user.Nickname,
		&user.HeaderImageURL,
		&user.ImageURL,
		&user.Mail,
		&user.Phone,
		&user.CardNumber,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s: %w", errs.ErrDB, methodName, err)
	}

	return &user, nil
}

func (r *ProfileRepository) CakesByUserID(ctx context.Context, userID uuid.UUID) ([]cakeDto.PreviewCakeDB, error) {
	methodName := "[Repo.CakesByUserID]"

	rows, err := r.db.QueryContext(ctx, querySelectCakesByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", errs.ErrDB, methodName, err)
	}

	defer rows.Close()
	var cakes []cakeDto.PreviewCakeDB
	for rows.Next() {
		var previewCake cakeDto.PreviewCakeDB
		if err = rows.Scan(
			&previewCake.ID,
			&previewCake.Name,
			&previewCake.PreviewImageURL,
			&previewCake.KgPrice,
			&previewCake.Rating,
			&previewCake.Description,
			&previewCake.Mass,
			&previewCake.DiscountKgPrice,
			&previewCake.DiscountEndTime,
			&previewCake.DateCreation,
			&previewCake.IsOpenForSale,
			&previewCake.OwnerID,
		); err != nil {
			return nil, fmt.Errorf("%w: %s: %w", errs.ErrDB, methodName, err)
		}

		cakes = append(cakes, previewCake)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %s: %w", errs.ErrDB, methodName, err)
	}

	return cakes, nil
}
