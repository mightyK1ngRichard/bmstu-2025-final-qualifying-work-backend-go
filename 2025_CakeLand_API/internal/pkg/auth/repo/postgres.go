package repo

import (
	"2025_CakeLand_API/internal/models/errs"
	"2025_CakeLand_API/internal/pkg/auth/dto"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/pkg/errors"
)

const (
	isUserExistsCommand            = `SELECT EXISTS(SELECT 1 FROM "user" WHERE mail = $1);`
	createUserCommand              = `INSERT INTO "user" (id, nickname, mail, password_hash, refresh_tokens_map) VALUES ($1, $2, $3, $4, $5);`
	getUserByEmailCommand          = `SELECT id, mail, refresh_tokens_map, password_hash FROM "user" WHERE mail = $1;`
	updateUserRefreshTokensCommand = `UPDATE "user" SET refresh_tokens_map = $1 WHERE id = $2;`
	getUserRefreshTokensCommand    = `SELECT refresh_tokens_map FROM "user" where id = $1`
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, in dto.CreateUserReq) error {
	methodName := "[Repo.CreateUser]"

	// Проверка существования пользователя с таким email
	var exists bool
	if err := r.db.QueryRowContext(ctx, isUserExistsCommand, in.Email).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrNotFound
		}
		return errors.Wrapf(err, "%v: %s", errs.ErrDB, methodName)
	}
	if exists {
		return errs.ErrAlreadyExists
	}

	// Сериализация RefreshTokensMap в JSON
	refreshTokensJSON, err := json.Marshal(in.RefreshTokensMap)
	if err != nil {
		return errors.Wrapf(err, "%s: failed to marshal refreshTokensJSON", methodName)
	}

	// Выполнение команды создания пользователя
	if _, err = r.db.ExecContext(ctx,
		createUserCommand,
		in.UUID,
		in.UUID, // изначально username пользователя равен его id
		in.Email,
		in.PasswordHash,
		refreshTokensJSON,
	); err != nil {
		return errors.Wrapf(err, "%v: %s", errs.ErrDB, methodName)
	}

	return nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, in dto.GetUserByEmailReq) (*dto.GetUserByEmailRes, error) {
	methodName := "[Repo.GetUserByEmail]"

	row := r.db.QueryRowContext(ctx, getUserByEmailCommand, in.Email)
	var res dto.GetUserByEmailRes
	if err := row.Scan(&res.ID, &res.Email, &res.RefreshTokensMap, &res.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, errors.Wrapf(err, "%v: %s", errs.ErrDB, methodName)
	}

	return &res, nil
}

func (r *AuthRepository) UpdateUserRefreshTokens(ctx context.Context, in dto.UpdateUserRefreshTokensReq) error {
	methodName := "[Repo.UpdateUserRefreshTokens]"

	// Сериализация RefreshTokensMap в JSON
	refreshTokensJSON, err := json.Marshal(in.RefreshTokensMap)
	if err != nil {
		return errors.Wrapf(err, "%s: failed to marshal refreshTokensJSON", methodName)
	}

	// Выполнение команды обновления токенов
	if _, err = r.db.ExecContext(ctx, updateUserRefreshTokensCommand, refreshTokensJSON, in.UserID); err != nil {
		return errors.Wrapf(err, "%v: %s", errs.ErrDB, methodName)
	}

	return nil
}

func (r *AuthRepository) GetUserRefreshTokens(ctx context.Context, in dto.GetUserRefreshTokensReq) (*dto.GetUserRefreshTokensRes, error) {
	methodName := "[Repo.GetUserRefreshTokens]"

	var refreshTokens []byte
	row := r.db.QueryRowContext(ctx, getUserRefreshTokensCommand, in.UserID)
	if err := row.Scan(&refreshTokens); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, errors.Wrapf(err, "%v: %s", errs.ErrDB, methodName)
	}

	var refreshTokensMap map[string]string
	if err := json.Unmarshal(refreshTokens, &refreshTokensMap); err != nil {
		return nil, errors.Wrapf(err, "%s: failed to unmarshal refreshTokensMap", methodName)
	}

	return &dto.GetUserRefreshTokensRes{
		RefreshTokensMap: refreshTokensMap,
	}, nil
}
