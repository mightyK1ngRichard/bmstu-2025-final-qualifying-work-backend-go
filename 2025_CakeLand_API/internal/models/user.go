package models

import (
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type RefreshTokenMap map[string]string // key: fingerprint, value: refreshToken

type User struct {
	ID               uuid.UUID       // Код
	FIO              null.String     // ФИО
	Address          null.String     // Адрес
	Nickname         string          // Уникальный псевдоним (default: id)
	ImageURL         null.String     // Картинка
	Mail             string          // Почта
	PasswordHash     []byte          // Пароль
	Phone            null.String     // Телефон
	CardNumber       null.String     // Номер кредитной карты
	RefreshTokensMap RefreshTokenMap // Рефреш токены (key: fingerprint, value: refreshToken)
}

func (u *User) ConvertToUserGRPC() *generated.User {
	var fio *wrapperspb.StringValue
	if u.FIO.Valid {
		fio = wrapperspb.String(u.FIO.String)
	}

	return &generated.User{
		Id:       u.ID.String(),
		Nickname: u.Nickname,
		Mail:     u.Mail,
		Fio:      fio,
	}
}

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
