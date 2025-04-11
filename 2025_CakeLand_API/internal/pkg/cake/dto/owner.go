package dto

import (
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Owner struct {
	ID             uuid.UUID
	FIO            null.String
	Address        null.String
	Nickname       string
	Mail           string
	Phone          null.String
	ImageURL       null.String
	HeaderImageURL null.String
}

func (pc *Owner) ConvertToGrpcUser() *generated.User {
	var fio *wrapperspb.StringValue
	if pc.FIO.Valid {
		fio = wrapperspb.String(pc.FIO.String)
	}

	var address *wrapperspb.StringValue
	if pc.Address.Valid {
		address = wrapperspb.String(pc.Address.String)
	}

	var phone *wrapperspb.StringValue
	if pc.Phone.Valid {
		phone = wrapperspb.String(pc.Phone.String)
	}

	var imageURL *wrapperspb.StringValue
	if pc.ImageURL.Valid {
		imageURL = wrapperspb.String(pc.ImageURL.String)
	}

	var headerImageURL *wrapperspb.StringValue
	if pc.HeaderImageURL.Valid {
		headerImageURL = wrapperspb.String(pc.HeaderImageURL.String)
	}

	return &generated.User{
		Id:             pc.ID.String(),
		Fio:            fio,
		Nickname:       pc.Nickname,
		Mail:           pc.Mail,
		Address:        address,
		Phone:          phone,
		ImageURL:       imageURL,
		HeaderImageURL: headerImageURL,
	}
}
