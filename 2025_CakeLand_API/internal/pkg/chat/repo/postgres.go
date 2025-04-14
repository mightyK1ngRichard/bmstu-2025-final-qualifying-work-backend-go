package repo

import (
	"2025_CakeLand_API/internal/models"
	"context"
	"database/sql"
)

const (
	queryAddMessage = "INSERT INTO message (id, text, date_creation, owner_id, receiver_id) VALUES ($1, $2, $3, $4, $5)"
)

type IChatRepository interface {
	AddMessage(context.Context, models.Message) error
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
	_, err := r.db.ExecContext(ctx, queryAddMessage, msg.ID, msg.Text, msg.DateCreation, msg.OwnerID, msg.ReceiverID)
	return err
}
