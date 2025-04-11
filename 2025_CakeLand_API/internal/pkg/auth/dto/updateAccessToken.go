package dto

import "time"

type UpdateAccessTokenReq struct {
	RefreshToken string
	Fingerprint  string
}

type UpdateAccessTokenRes struct {
	AccessToken string
	ExpiresIn   time.Time
}
