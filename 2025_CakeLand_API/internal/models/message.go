package models

import "time"

type Message struct {
	ID           string
	Text         string
	OwnerID      string
	ReceiverID   string
	DateCreation time.Time
}
