package repo

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	queryAddMessage        = "INSERT INTO message (id, text, date_creation, owner_id, receiver_id) VALUES ($1, $2, $3, $4, $5)"
	queryUserInterlocutors = `SELECT receiver_id FROM message WHERE owner_id = $1`
	queryUserByID          = `
		SELECT id,
			   fio,
			   address,
			   nickname,
			   image_url,
			   mail,
			   phone,
			   header_image_url
		FROM "user"
		WHERE id = $1;
	`
	queryUserHistory = `
		SELECT id, text, date_creation, owner_id, receiver_id 
		FROM message 
		WHERE (owner_id = $1 AND receiver_id = $2) OR (owner_id = $2 AND receiver_id = $1)
	`
)

type IChatRepository interface {
	AddMessage(context.Context, models.Message) error
	UserInterlocutors(context.Context, string) ([]string, error)
	UserByID(context.Context, string) (*models.User, error)
	ChatHistory(ctx context.Context, ownerID, interlocutorID string) ([]*models.Message, error)
}

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (r *ChatRepository) AddMessage(ctx context.Context, msg models.Message) error {
	methodName := "[Repo.AddMessage]"
	_, err := r.db.ExecContext(ctx, queryAddMessage, msg.ID, msg.Text, msg.DateCreation, msg.OwnerID, msg.ReceiverID)
	if err != nil {
		return errs.WrapDBError(methodName, err)
	}

	return nil
}

func (r *ChatRepository) UserByID(ctx context.Context, userID string) (*models.User, error) {
	methodName := "[Repo.GetUserByID]"

	var user models.User
	if err := r.db.QueryRowContext(ctx, queryUserByID, userID).Scan(
		&user.ID,
		&user.FIO,
		&user.Address,
		&user.Nickname,
		&user.ImageURL,
		&user.Mail,
		&user.Phone,
		&user.HeaderImageURL,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.WrapDBError(methodName, err)
	}

	return &user, nil
}

func (r *ChatRepository) ChatHistory(ctx context.Context, ownerID, interlocutorID string) ([]*models.Message, error) {
	methodName := "[Repo.GetChatHistory]"

	rows, err := r.db.QueryContext(ctx, queryUserHistory, ownerID, interlocutorID)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	defer rows.Close()
	var messages []*models.Message
	for rows.Next() {
		var message models.Message
		if err = rows.Scan(
			&message.ID,
			&message.Text,
			&message.DateCreation,
			&message.OwnerID,
			&message.ReceiverID,
		); err != nil {
			return nil, fmt.Errorf("%w: %s", errs.ErrDB, methodName)
		}
		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return messages, nil
}

func (r *ChatRepository) UserInterlocutors(ctx context.Context, userID string) ([]string, error) {
	methodName := "[Repo.GetUserInterlocutors]"

	rows, err := r.db.QueryContext(ctx, queryUserInterlocutors, userID)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, errs.WrapDBError(methodName, err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
