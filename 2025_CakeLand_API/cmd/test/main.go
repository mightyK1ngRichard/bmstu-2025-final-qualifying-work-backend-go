package main

import (
	"2025_CakeLand_API/internal/pkg/config"
	"2025_CakeLand_API/internal/pkg/profile/repo"
	"2025_CakeLand_API/internal/pkg/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
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

	ctx := context.Background()
	rep := repo.NewProfileRepository(db)
	userIDStr := "550e8400-e29b-41d4-a716-446655440021"
	userID, _ := uuid.Parse(userIDStr)

	data, err := rep.CakesByUserID(ctx, userID)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(data.ID, data.Nickname, data.ImageURL)
	fmt.Println(data)
}
