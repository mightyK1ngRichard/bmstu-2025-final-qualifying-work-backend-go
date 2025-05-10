package repo

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

const (
	queryAddNotification = `
		INSERT INTO notification (
			id,
			title,
			content,
			date_creation,
			sender_id,
			recipient_id,
			notification_kind,
		    order_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	queryGetUserNotifications = `
		SELECT id, title, content, date_creation, sender_id, recipient_id, notification_kind, order_id
		FROM notification
		WHERE recipient_id = $1
		ORDER BY date_creation DESC;
	`
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) CreateNotification(ctx context.Context, in models.Notification) error {
	const methodName = "[NotificationRepository.CreateNotification]"

	if _, err := r.db.ExecContext(ctx, queryAddNotification,
		in.ID,
		in.Title,
		in.Message,
		in.CreatedAt,
		in.SenderID,
		in.RecipientID,
		string(in.NotificationKind),
		in.OrderID,
	); err != nil {
		return errors.Wrap(err, methodName)
	}

	return nil
}

func (r *NotificationRepository) GetNotifications(ctx context.Context, userID string) ([]models.Notification, error) {
	const methodName = "[NotificationRepository.GetNotifications]"

	rows, err := r.db.QueryContext(ctx, queryGetUserNotifications, userID)
	if err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		if err = rows.Scan(
			&notification.ID,
			&notification.Title,
			&notification.Message,
			&notification.CreatedAt,
			&notification.SenderID,
			&notification.RecipientID,
			&notification.NotificationKind,
			&notification.OrderID,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errs.WrapDBError(methodName, errs.ErrNotFound)
			}
			return nil, errs.WrapDBError(methodName, err)
		}

		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapDBError(methodName, err)
	}

	return notifications, nil
}
