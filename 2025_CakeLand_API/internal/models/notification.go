package models

import (
	gen "2025_CakeLand_API/internal/pkg/notification/delivery/grpc/generated"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type NotificationKind string

const (
	KindMessage     NotificationKind = "message"      // Личное сообщение
	KindFeedback    NotificationKind = "feedback"     // Отзыв
	KindOrderUpdate NotificationKind = "order_update" // Обновление по заказу
	KindSystem      NotificationKind = "system"       // Системное уведомление
	KindPromo       NotificationKind = "promo"        // Рекламное уведомление
)

type Notification struct {
	ID               uuid.UUID
	Title            string
	Message          string
	CreatedAt        time.Time
	SenderID         string
	RecipientID      string
	CakeID           null.String
	NotificationKind NotificationKind
}

func ConvertProtoNotificationKind(protoKind gen.NotificationKind) NotificationKind {
	switch protoKind {
	case gen.NotificationKind_MESSAGE:
		return KindMessage
	case gen.NotificationKind_FEEDBACK:
		return KindFeedback
	case gen.NotificationKind_ORDER_UPDATE:
		return KindOrderUpdate
	case gen.NotificationKind_SYSTEM:
		return KindSystem
	case gen.NotificationKind_PROMO:
		return KindPromo
	default:
		return ""
	}
}

func convertNotificationKindToProto(kind NotificationKind) gen.NotificationKind {
	switch kind {
	case KindMessage:
		return gen.NotificationKind_MESSAGE
	case KindFeedback:
		return gen.NotificationKind_FEEDBACK
	case KindOrderUpdate:
		return gen.NotificationKind_ORDER_UPDATE
	case KindSystem:
		return gen.NotificationKind_SYSTEM
	case KindPromo:
		return gen.NotificationKind_PROMO
	default:
		return gen.NotificationKind(0)
	}
}

func (n *Notification) ToProto() *gen.Notification {
	createdAtTimestamp := timestamppb.New(n.CreatedAt)

	protoNotif := &gen.Notification{
		Id:        n.ID.String(),
		Title:     n.Title,
		Message:   n.Message,
		CreatedAt: createdAtTimestamp,
		SenderID:  n.SenderID,
		Kind:      convertNotificationKindToProto(n.NotificationKind),
	}

	if n.CakeID.Valid {
		protoNotif.CakeID = &n.CakeID.String
	}

	return protoNotif
}
