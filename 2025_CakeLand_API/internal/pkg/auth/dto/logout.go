package dto

type LogoutReq struct {
	RefreshToken string
	Fingerprint  string
}

type LogoutRes struct {
	Message string
}
