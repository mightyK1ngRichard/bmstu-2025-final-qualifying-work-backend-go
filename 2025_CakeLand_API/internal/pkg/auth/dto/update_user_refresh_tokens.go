package dto

type UpdateUserRefreshTokensReq struct {
	UserID           string
	RefreshTokensMap map[string]string
}
