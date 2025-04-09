package dto

import (
	"2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type RefreshTokenMap map[string]string // key: fingerprint, value: refreshToken

// Scan Реализуем интерфейс sql.Scanner для JSONB
func (rtm *RefreshTokenMap) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("RefreshTokenMap: type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, rtm)
}

type Profile struct {
	ID               uuid.UUID
	FIO              null.String
	Address          null.String
	Nickname         string
	ImageURL         null.String
	HeaderImageURL   null.String
	Mail             string
	PasswordHash     string
	Phone            null.String
	CardNumber       null.String
	RefreshTokensMap RefreshTokenMap
}

func (p *Profile) ConvertToGrpcModel() *generated.Profile {
	return &generated.Profile{
		Id:               p.ID.String(),
		Fio:              stringOrNil(p.FIO.NullString),
		Address:          stringOrNil(p.Address.NullString),
		Nickname:         p.Nickname,
		ImageUrl:         stringOrNil(p.ImageURL.NullString),
		HeaderImageUrl:   stringOrNil(p.HeaderImageURL.NullString),
		Mail:             p.Mail,
		PasswordHash:     p.PasswordHash,
		Phone:            stringOrNil(p.Phone.NullString),
		CardNumber:       stringOrNil(p.CardNumber.NullString),
		RefreshTokensMap: p.RefreshTokensMap,
	}
}

// Вспомогательная функция для nullable строк
func stringOrNil(s sql.NullString) *wrapperspb.StringValue {
	if s.Valid {
		return wrapperspb.String(s.String)
	}

	return nil
}
