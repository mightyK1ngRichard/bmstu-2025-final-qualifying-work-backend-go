package repo

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	cakeDto "2025_CakeLand_API/internal/pkg/cake/dto"
	gen "2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/profile/dto"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

const (
	querySelectProfileByID = `
		SELECT id, fio, nickname, header_image_url, image_url, mail, phone, card_number FROM "user" WHERE id = $1 LIMIT 1
	`
	querySelectCakesByUserID = `
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
			   owner_id,
			   model_3d_url
		FROM cake
		WHERE owner_id = $1;
    `
	queryCreateAddress = `
		INSERT INTO address (id, user_id, latitude, longitude, formatted_address) VALUES ($1, $2, $3, $4, $5)
	`
	queryGetUserAddresses = `
		SELECT id,
			   user_id,
			   latitude,
			   longitude,
			   formatted_address,
			   entrance,
			   floor,
			   apartment,
			   comment
		FROM address
		WHERE user_id = $1
	`
	queryUpdateUserAddress = `
		UPDATE address
		SET 
			entrance = COALESCE($1, entrance),
			floor = COALESCE($2, floor),
			apartment = COALESCE($3, apartment),
			comment = COALESCE($4, comment),
			updated_at = now()
		WHERE id = $5 AND user_id = $6
		RETURNING id, user_id, latitude, longitude, formatted_address, entrance, floor, apartment, comment
	`
	queryUpdateUserAvatar = `UPDATE "user" SET image_url = $1 WHERE id = $2`
	queryUpdateUserHeader = `UPDATE "user" SET header_image_url = $1 WHERE id = $2`
	queryUpdateUserData   = `UPDATE "user" SET nickname = $1, fio = $2 WHERE id = $3`
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) UpdateUserAddresses(ctx context.Context, userID uuid.UUID, req *gen.UpdateUserAddressesReq) (models.Address, error) {
	const methodName = "[ProfileRepository.UpdateUserAddresses]"

	row := r.db.QueryRowContext(ctx, queryUpdateUserAddress,
		req.Entrance,
		req.Floor,
		req.Apartment,
		req.Comment,
		req.AddressID,
		userID,
	)

	var addr models.Address
	if err := row.Scan(
		&addr.ID,
		&addr.UserID,
		&addr.Latitude,
		&addr.Longitude,
		&addr.FormattedAddress,
		&addr.Entrance,
		&addr.Floor,
		&addr.Apartment,
		&addr.Comment,
	); err != nil {
		return models.Address{}, errs.WrapDBError(methodName, err)
	}

	return addr, nil
}

func (r *ProfileRepository) GetUserAddresses(ctx context.Context, userID uuid.UUID) ([]models.Address, error) {
	const methodName = "[ProfileRepository.GetUserAddresses]"

	rows, err := r.db.QueryContext(ctx, queryGetUserAddresses, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.Address
	for rows.Next() {
		var addr models.Address
		if err = rows.Scan(
			&addr.ID,
			&addr.UserID,
			&addr.Latitude,
			&addr.Longitude,
			&addr.FormattedAddress,
			&addr.Entrance,
			&addr.Floor,
			&addr.Apartment,
			&addr.Comment,
		); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}

		addresses = append(addresses, addr)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return addresses, nil
}

func (r *ProfileRepository) CreateAddress(ctx context.Context, address *models.Address) error {
	const methodName = "[ProfileRepository.CreateAddress]"

	if _, err := r.db.ExecContext(ctx, queryCreateAddress,
		address.ID,
		address.UserID,
		address.Latitude,
		address.Longitude,
		address.FormattedAddress,
	); err != nil {
		return errs.WrapDBError(methodName, err)
	}

	return nil
}

func (r *ProfileRepository) UserInfo(ctx context.Context, userID uuid.UUID) (*dto.Profile, error) {
	const methodName = "[ProfileRepository.UserInfo]"

	var user dto.Profile
	if err := r.db.QueryRowContext(ctx, querySelectProfileByID, userID).Scan(
		&user.ID,
		&user.FIO,
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
		return nil, errs.WrapDBError(methodName, err)
	}

	return &user, nil
}

func (r *ProfileRepository) CakesByUserID(ctx context.Context, userID uuid.UUID) ([]cakeDto.PreviewCakeDB, error) {
	const methodName = "[ProfileRepository.CakesByUserID]"

	rows, err := r.db.QueryContext(ctx, querySelectCakesByUserID, userID)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
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
			&previewCake.ReviewsCount,
			&previewCake.StarsSum,
			&previewCake.Description,
			&previewCake.Mass,
			&previewCake.DiscountKgPrice,
			&previewCake.DiscountEndTime,
			&previewCake.DateCreation,
			&previewCake.IsOpenForSale,
			&previewCake.OwnerID,
			&previewCake.Model3DURL,
		); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}

		cakes = append(cakes, previewCake)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return cakes, nil
}

func (r *ProfileRepository) UpdateUserAvatar(ctx context.Context, userID uuid.UUID, imageURL string) error {
	const methodName = "[ProfileRepository.UpdateUserAvatar]"

	if _, err := r.db.ExecContext(ctx, queryUpdateUserAvatar, imageURL, userID); err != nil {
		return errs.WrapDBError(methodName, err)
	}

	return nil
}

func (r *ProfileRepository) UpdateUserHeaderImage(ctx context.Context, userID uuid.UUID, imageURL string) error {
	const methodName = "[ProfileRepository.UpdateUserHeaderImage]"

	if _, err := r.db.ExecContext(ctx, queryUpdateUserHeader, imageURL, userID); err != nil {
		return errs.WrapDBError(methodName, err)
	}

	return nil
}

func (r *ProfileRepository) UpdateUserData(ctx context.Context, userID uuid.UUID, in *gen.UpdateUserDataReq) error {
	const methodName = "[ProfileRepository.UpdateUserData]"

	if _, err := r.db.ExecContext(ctx, queryUpdateUserData,
		in.UpdatedUserName,
		in.UpdatedFIO,
		userID,
	); err != nil {
		return errs.WrapDBError(methodName, err)
	}

	return nil
}
