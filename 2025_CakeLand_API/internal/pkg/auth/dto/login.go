package dto

import "time"

type LoginReq struct {
	Email       string
	Password    string
	Fingerprint string
}

type LoginRes struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Time
}
