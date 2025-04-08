package main

import (
	"2025_CakeLand_API/internal/models"
	"2025_CakeLand_API/internal/pkg/cake/repo"
	"2025_CakeLand_API/internal/pkg/config"
	"2025_CakeLand_API/internal/pkg/utils"
	"context"
	"fmt"
	"log"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := utils.ConnectPostgres(&conf.DB)
	if err != nil {
		log.Fatal(err)
	}

	rep := repo.NewCakeRepository(db)
	categories, err := rep.CategoryIDsByGenderName(context.Background(), models.GenderFemale)
	for _, category := range categories {
		fmt.Println(category)
	}
}
