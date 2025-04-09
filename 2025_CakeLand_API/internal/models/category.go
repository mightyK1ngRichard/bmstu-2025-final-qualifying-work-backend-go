package models

import (
	"2025_CakeLand_API/internal/pkg/cake/delivery/grpc/generated"
	"github.com/google/uuid"
)

type CategoryGender string

const (
	GenderMale        CategoryGender = "male"
	GenderFemale      CategoryGender = "female"
	GenderChild       CategoryGender = "child"
	GenderUnspecified CategoryGender = "unspecified"
)

func ConvertToCategoryGender(categoryGenStr string) CategoryGender {
	switch categoryGenStr {
	case string(GenderMale):
		return GenderMale
	case string(GenderFemale):
		return GenderFemale
	case string(GenderChild):
		return GenderChild
	default:
		return GenderUnspecified
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

// Category Модель категории
type Category struct {
	ID              uuid.UUID        // Код
	Name            string           // Название
	ImageURL        string           // Картинка
	CategoryGenders []CategoryGender // Теги торта
}

func (c *Category) ConvertToCategoryGRPC() *generated.Category {
	// Создаем пустой слайс для gender_tags
	var genderTags []generated.CategoryGender

	// Заполняем слайс genderTags значениями из CategoryGenders
	for _, gender := range c.CategoryGenders {
		genderTags = append(genderTags, gender.ConvertToGRPCCategoryGender())
	}

	return &generated.Category{
		Id:         c.ID.String(),
		Name:       c.Name,
		ImageUrl:   c.ImageURL,
		GenderTags: genderTags,
	}
}
