package models

import (
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	profileGen "2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
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
	Nickname         string          // Уникальный псевдоним (default: id)
	ImageURL         null.String     // Картинка
	HeaderImageURL   null.String     // Картинка шапки профиля
	Mail             string          // Почта
	PasswordHash     []byte          // Пароль
	Phone            null.String     // Телефон
	CardNumber       null.String     // Номер кредитной карты
	RefreshTokensMap RefreshTokenMap // Рефреш токены (key: fingerprint, value: refreshToken)
}

type UserInfo struct {
	ID             string
	FIO            null.String
	Address        null.String
	Nickname       string
	ImageURL       null.String
	HeaderImageURL null.String
	Mail           string
	Phone          null.String
}

func NewUserInfo(u *profileGen.Profile) *UserInfo {
	return &UserInfo{
		ID:       u.GetId(),
		Nickname: u.GetNickname(),
		Mail:     u.GetMail(),
		FIO: null.NewString(
			u.GetFio().GetValue(),
			u.Fio != nil,
		),
		Phone: null.NewString(
			u.GetPhone().GetValue(),
			u.Phone != nil,
		),
		ImageURL: null.NewString(
			u.GetImageUrl().GetValue(),
			u.ImageUrl != nil,
		),
		HeaderImageURL: null.NewString(
			u.GetHeaderImageUrl().GetValue(),
			u.HeaderImageUrl != nil,
		),
	}
}

func (u *User) ConvertToUserGRPC() *generated.User {
	var fio *wrapperspb.StringValue
	if u.FIO.Valid {
		fio = wrapperspb.String(u.FIO.String)
	}

	var phoneNumber *wrapperspb.StringValue
	if u.Phone.Valid {
		phoneNumber = wrapperspb.String(u.Phone.String)
	}

	var imageURL *wrapperspb.StringValue
	if u.ImageURL.Valid {
		imageURL = wrapperspb.String(u.ImageURL.String)
	}

	var headerImageURL *wrapperspb.StringValue
	if u.HeaderImageURL.Valid {
		headerImageURL = wrapperspb.String(u.HeaderImageURL.String)
	}

	return &generated.User{
		Id:             u.ID.String(),
		Nickname:       u.Nickname,
		Mail:           u.Mail,
		Fio:            fio,
		Phone:          phoneNumber,
		ImageURL:       imageURL,
		HeaderImageURL: headerImageURL,
	}
}

func (u *UserInfo) ConvertToGRPCProfile() *profileGen.Profile {
	var fio *wrapperspb.StringValue
	if u.FIO.Valid {
		fio = wrapperspb.String(u.FIO.String)
	}

	var phoneNumber *wrapperspb.StringValue
	if u.Phone.Valid {
		phoneNumber = wrapperspb.String(u.Phone.String)
	}

	var imageURL *wrapperspb.StringValue
	if u.ImageURL.Valid {
		imageURL = wrapperspb.String(u.ImageURL.String)
	}

	var headerImageURL *wrapperspb.StringValue
	if u.HeaderImageURL.Valid {
		headerImageURL = wrapperspb.String(u.HeaderImageURL.String)
	}

	return &profileGen.Profile{
		Id:             u.ID,
		Nickname:       u.Nickname,
		Mail:           u.Mail,
		Fio:            fio,
		Phone:          phoneNumber,
		ImageUrl:       imageURL,
		HeaderImageUrl: headerImageURL,
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
