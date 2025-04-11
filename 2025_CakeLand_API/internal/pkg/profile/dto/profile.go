package dto

import (
	"2025_CakeLand_API/internal/pkg/cake/dto"
	"2025_CakeLand_API/internal/pkg/profile/delivery/grpc/generated"
	"database/sql"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Profile struct {
	ID             uuid.UUID
	FIO            null.String
	Address        null.String
	Nickname       string
	ImageURL       null.String
	HeaderImageURL null.String
	Mail           string
	Phone          null.String
	CardNumber     null.String
}

func (p *Profile) ConvertToGrpcModel() *generated.Profile {
	return &generated.Profile{
		Id:             p.ID.String(),
		Fio:            stringOrNil(p.FIO.NullString),
		Address:        stringOrNil(p.Address.NullString),
		Nickname:       p.Nickname,
		ImageUrl:       stringOrNil(p.ImageURL.NullString),
		HeaderImageUrl: stringOrNil(p.HeaderImageURL.NullString),
		Mail:           p.Mail,
		Phone:          stringOrNil(p.Phone.NullString),
		CardNumber:     stringOrNil(p.CardNumber.NullString),
	}
}

func (p *Profile) ConvertToOwner() dto.Owner {
	return dto.Owner{
		ID:             p.ID,
		Nickname:       p.Nickname,
		Mail:           p.Mail,
		FIO:            p.FIO,
		Address:        p.Address,
		HeaderImageURL: p.HeaderImageURL,
		ImageURL:       p.ImageURL,
		Phone:          p.Phone,
	}
}

// Вспомогательная функция для nullable строк
func stringOrNil(s sql.NullString) *wrapperspb.StringValue {
	if s.Valid {
		return wrapperspb.String(s.String)
	}

	return nil
}
