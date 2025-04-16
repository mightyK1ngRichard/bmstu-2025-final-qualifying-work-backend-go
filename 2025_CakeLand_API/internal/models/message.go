package models

import (
	gen "2025_CakeLand_API/internal/pkg/chat/delivery/grpc/generated"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Message struct {
	ID           string
	Text         string
	OwnerID      string
	ReceiverID   string
	DateCreation time.Time
}

func (m *Message) ConvertToGrpcModel() *gen.ChatMessage {
	return &gen.ChatMessage{
		Id:             m.ID,
		InterlocutorID: m.ReceiverID,
		SenderID:       m.OwnerID,
		Text:           m.Text,
		DateCreation:   timestamppb.New(m.DateCreation),
	}
}
