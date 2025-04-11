package dto

type GetUserRefreshTokensReq struct {
	UserID string
}

type GetUserRefreshTokensRes struct {
	RefreshTokensMap map[string]string
}
