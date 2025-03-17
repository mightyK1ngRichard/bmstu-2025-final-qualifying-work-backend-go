package models

import "time"

type JWTTokenPayload struct {
	UserUID   string
	Token     string
	ExpiresIn time.Time
}
