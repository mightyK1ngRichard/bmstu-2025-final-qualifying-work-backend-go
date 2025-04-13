package models

import (
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CategoryGender string

const (
	GenderMale        CategoryGender = "male"
	GenderFemale      CategoryGender = "female"
	GenderChild       CategoryGender = "child"
	GenderUnspecified CategoryGender = "unspecified"
)

func (c CategoryGender) ConvertToGRPCCategoryGender() generated.CategoryGender {
	switch c {
	case GenderMale:
		return generated.CategoryGender_MALE
	case GenderFemale:
		return generated.CategoryGender_FEMALE
	case GenderChild:
		return generated.CategoryGender_CHILD
	default:
		return generated.CategoryGender_CATEGORY_GENDER_UNSPECIFIED
	}
}

func ConvertToCategoryGenderFromGrpc(categoryGen generated.CategoryGender) CategoryGender {
	switch categoryGen {
	case generated.CategoryGender_MALE:
		return GenderMale
	case generated.CategoryGender_FEMALE:
		return GenderFemale
	case generated.CategoryGender_CHILD:
		return GenderChild
	default:
		return GenderUnspecified
	}
}

func ParseGenderTags(tags pq.StringArray) []CategoryGender {
	genders := make([]CategoryGender, len(tags))
	for i, tag := range tags {
		switch tag {
		case string(GenderMale):
			genders[i] = GenderMale
		case string(GenderFemale):
			genders[i] = GenderFemale
		case string(GenderChild):
			genders[i] = GenderChild
		case string(GenderUnspecified):
			genders[i] = GenderUnspecified
		}
	}
	return genders
}

// Category Модель категории
type Category struct {
	ID              uuid.UUID        // Код
	Name            string           // Название
	ImageURL        string           // Картинка
	CategoryGenders []CategoryGender // Теги торта
}

func (c *Category) ConvertToCategoryGRPC() *generated.Category {
	// Создаем пустой слайс для gender_tags
	genderTags := make([]generated.CategoryGender, len(c.CategoryGenders))

	// Заполняем слайс genderTags значениями из CategoryGenders
	for i, gender := range c.CategoryGenders {
		genderTags[i] = gender.ConvertToGRPCCategoryGender()
	}

	return &generated.Category{
		Id:         c.ID.String(),
		Name:       c.Name,
		ImageUrl:   c.ImageURL,
		GenderTags: genderTags,
	}
}
