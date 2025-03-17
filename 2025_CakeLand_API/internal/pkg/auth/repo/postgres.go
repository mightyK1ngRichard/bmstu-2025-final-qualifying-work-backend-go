package repo

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/auth/entities"
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

func (r *AuthRepository) CreateUser(ctx context.Context, in entities.CreateUserReq) error {
	// Проверка существования пользователя с таким email
	var exists bool
	err := r.db.QueryRowContext(ctx, isUserExistsCommand, in.Email).Scan(&exists)
	if err != nil {
		return errors.Wrap(err, "ошибка при проверке существования пользователя")
	}
	if exists {
		return models.ErrUserAlreadyExists
	}

	// Сериализация RefreshTokensMap в JSON
	refreshTokensJSON, err := json.Marshal(in.RefreshTokensMap)
	if err != nil {
		return errors.Wrap(err, "ошибка сериализации RefreshTokensMap в JSON при создании пользователя")
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
		return errors.Wrap(err, "ошибка выполнения команды создания пользователя в базе данных")
	}

	return nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, in entities.GetUserByEmailReq) (*entities.GetUserByEmailRes, error) {
	row := r.db.QueryRowContext(ctx, getUserByEmailCommand, in.Email)
	var res entities.GetUserByEmailRes
	var refreshTokensMap []byte
	if err := row.Scan(&res.ID, &res.Email, &refreshTokensMap, &res.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, errors.Wrap(err, "ошибка получения данных пользователя из базы данных")
	}
	if err := json.Unmarshal(refreshTokensMap, &res.RefreshTokensMap); err != nil {
		return nil, errors.Wrapf(err, "ошибка декодирования JSON refreshTokensMap для пользователя с email %s", in.Email)
	}

	return &res, nil
}

func (r *AuthRepository) UpdateUserRefreshTokens(ctx context.Context, in entities.UpdateUserRefreshTokensReq) error {
	// Сериализация RefreshTokensMap в JSON
	refreshTokensJSON, err := json.Marshal(in.RefreshTokensMap)
	if err != nil {
		return errors.Wrap(err, "ошибка сериализации RefreshTokensMap в JSON при обновлении токенов пользователя")
	}

	// Выполнение команды обновления токенов
	if _, err = r.db.ExecContext(ctx, updateUserRefreshTokensCommand, refreshTokensJSON, in.UserID); err != nil {
		return errors.Wrapf(err, "ошибка выполнения команды обновления токенов пользователя с ID %s", in.UserID)
	}

	return nil
}

func (r *AuthRepository) GetUserRefreshTokens(ctx context.Context, in entities.GetUserRefreshTokensReq) (*entities.GetUserRefreshTokensRes, error) {
	var refreshTokens []byte
	row := r.db.QueryRowContext(ctx, getUserRefreshTokensCommand, in.UserID)
	if err := row.Scan(&refreshTokens); err != nil {
		return nil, errors.New(`refresh_tokens_map not found`)
	}
	var refreshTokensMap map[string]string
	if err := json.Unmarshal(refreshTokens, &refreshTokensMap); err != nil {
		return nil, errors.Wrapf(err, `ошибка декодирования JSON refreshTokensMap`)
	}
	return &entities.GetUserRefreshTokensRes{
		RefreshTokensMap: refreshTokensMap,
	}, nil
}
