package handler

import (
	"2025_CakeLand_API/internal/domains"
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/models/errs"
	gen "2025_CakeLand_API/internal/pkg/notification/delivery/grpc/generated"
	"2025_CakeLand_API/internal/pkg/notification/repo"
	"2025_CakeLand_API/internal/pkg/utils/jwt"
	md "2025_CakeLand_API/internal/pkg/utils/metadata"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
	"sync"
	"time"
)

type NotificationStream struct {
	ch   chan *gen.NotificationResponse
	done chan struct{}
}

type NotificationHandler struct {
	gen.UnimplementedNotificationServiceServer

	log        *slog.Logger
	mdProvider *md.MetadataProvider
	tokenator  *jwt.Tokenator
	repo       *repo.NotificationRepository

	streamsMu sync.Mutex
	streams   map[string]*NotificationStream
}

func NewNotificationHandler(
	log *slog.Logger,
	mdProvider *md.MetadataProvider,
	repo *repo.NotificationRepository,
	tokenator *jwt.Tokenator,
) *NotificationHandler {
	return &NotificationHandler{
		mdProvider: mdProvider,
		log:        log,
		repo:       repo,
		tokenator:  tokenator,
		streams:    make(map[string]*NotificationStream),
	}
}

func (h *NotificationHandler) StreamNotifications(_ *emptypb.Empty, stream grpc.ServerStreamingServer[gen.NotificationResponse]) error {
	// Достаём ID пользователя
	userID, err := h.extractUserIDFromContext(stream.Context())
	if err != nil {
		return errs.ConvertToGrpcError(stream.Context(), h.log, err, "failed to extract user ID")
	}

	streamID := userID
	notifStream := &NotificationStream{
		ch:   make(chan *gen.NotificationResponse, 10),
		done: make(chan struct{}),
	}

	h.streamsMu.Lock()
	h.streams[streamID] = notifStream
	h.streamsMu.Unlock()

	defer func() {
		h.streamsMu.Lock()
		delete(h.streams, streamID)
		close(notifStream.done)
		close(notifStream.ch)
		h.streamsMu.Unlock()
	}()

	for {
		select {
		case <-notifStream.done:
			return nil
		case <-stream.Context().Done(): // клиент закрыл соединение
			h.log.Info("client disconnected", slog.String("userID", userID))
			return nil
		case msg := <-notifStream.ch:
			if err = stream.Send(msg); err != nil {
				return err
			}
		}
	}
}

func (h *NotificationHandler) GetNotifications(ctx context.Context, _ *emptypb.Empty) (*gen.GetNotificationsResponse, error) {
	// Достаём userID из метаданных
	userID, err := h.extractUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Бизнес логика
	notifications, err := h.repo.GetNotifications(ctx, userID)
	if err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "error getting notifications")
	}

	// Ответ
	response := make([]*gen.Notification, len(notifications))
	for i, notification := range notifications {
		response[i] = notification.ToProto()
	}
	return &gen.GetNotificationsResponse{
		Notifications: response,
	}, nil
}

func (h *NotificationHandler) CreateNotification(ctx context.Context, in *gen.CreateNotificationRequest) (*gen.NotificationResponse, error) {
	// Достаём userID из метаданных
	userID, err := h.extractUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Запись в БД
	notification := models.Notification{
		ID:               uuid.New(),
		Title:            in.Title,
		Message:          in.Message,
		CreatedAt:        time.Now(),
		RecipientID:      in.RecipientID,
		SenderID:         userID,
		NotificationKind: models.ConvertProtoNotificationKind(in.Kind),
		CakeID: func(s *string) null.String {
			if s != nil {
				return null.StringFrom(*s)
			}
			return null.String{}
		}(in.CakeID),
	}
	if err = h.repo.CreateNotification(ctx, notification); err != nil {
		return nil, errs.ConvertToGrpcError(ctx, h.log, err, "failed to create notification")
	}

	// Ответ
	response := &gen.NotificationResponse{
		Notification: notification.ToProto(),
	}

	go func() {
		h.streamsMu.Lock()
		defer h.streamsMu.Unlock()

		// Отправить получателю
		if stream, ok := h.streams[notification.RecipientID]; ok {
			select {
			case stream.ch <- response:
			default:
				// буфер переполнен — можно залогировать или пропустить
			}
		}

		// Отправить отправителю (если отличается от получателя)
		if notification.SenderID != notification.RecipientID {
			if stream, ok := h.streams[notification.SenderID]; ok {
				select {
				case stream.ch <- response:
				default:
					// буфер переполнен
				}
			}
		}
	}()

	return response, nil
}

func (h *NotificationHandler) extractUserIDFromContext(ctx context.Context) (string, error) {
	accessToken, err := h.mdProvider.GetValue(ctx, domains.KeyAuthorization)
	if err != nil {
		return "", errs.ConvertToGrpcError(ctx, h.log, err,
			fmt.Sprintf("missing required metadata: %s", domains.KeyAuthorization),
		)
	}

	userID, err := h.tokenator.GetUserIDFromToken(accessToken, false)
	if err != nil {
		return "", errs.ConvertToGrpcError(ctx, h.log, err,
			"failed to extract userID from token",
		)
	}

	return userID, nil
}
