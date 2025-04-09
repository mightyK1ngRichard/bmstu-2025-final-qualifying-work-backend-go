package repo

import (
	cakeDto "2025_CakeLand_API/internal/pkg/cake/dto"
	"2025_CakeLand_API/internal/pkg/profile/dto"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

const (
	querySelectProfileByID   = `SELECT id, fio, address, nickname, header_image_url, image_url, mail, password_hash, phone, card_number, refresh_tokens_map FROM "user" WHERE id = $1 LIMIT 1`
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
	var user dto.Profile
	if err := r.db.QueryRowContext(ctx, querySelectProfileByID, userID).Scan(
		&user.ID,
		&user.FIO,
		&user.Address,
		&user.Nickname,
		&user.HeaderImageURL,
		&user.ImageURL,
		&user.Mail,
		&user.PasswordHash,
		&user.Phone,
		&user.CardNumber,
		&user.RefreshTokensMap,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ProfileRepository) CakesByUserID(ctx context.Context, userID uuid.UUID) ([]cakeDto.PreviewCake, error) {
	rows, err := r.db.QueryContext(ctx, querySelectCakesByUserID, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var cakes []cakeDto.PreviewCake
	for rows.Next() {
		var previewCake cakeDto.PreviewCake
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
			return nil, err
		}

		cakes = append(cakes, previewCake)
	}

	return cakes, rows.Err()
}
