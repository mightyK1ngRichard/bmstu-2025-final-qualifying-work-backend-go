package dto

import "time"

type RegisterReq struct {
	Email       string
	Password    string
	Nickname    string
	Fingerprint string
}

type RegisterRes struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Time
}
