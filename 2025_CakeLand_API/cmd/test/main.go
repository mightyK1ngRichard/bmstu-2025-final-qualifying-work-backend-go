package main

import (
	"2025_CakeLand_API/internal/models/errs"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

func main() {
	//conf, err := config.NewConfig()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//db, err := utils.ConnectPostgres(&conf.DB)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//ctx := context.Background()
	//rep := repo.NewProfileRepository(db)
	//userIDStr := "550e8400-e29b-41d4-a716-446655440021"
	//userID, _ := uuid.Parse(userIDStr)
	//
	//data, err := rep.CakesByUserID(ctx, userID)
	//if err != nil {
	//	log.Fatal(err)
	//}
	////fmt.Println(data.ID, data.Nickname, data.ImageURL)
	//fmt.Println(data)

	//log := logger.NewLogger("local")
	//ctx := context.Background()
	//err := errs.ConvertToGrpcError(ctx, log, errs.ErrNoMetadata, "refreshToken is empty")
	//jsonErr := grpcErrorToJSON(err)
	//fmt.Println(jsonErr)

	err := errs.ErrNotFound
	fmt.Println(errors.Wrapf(err, "text"))
}

func grpcErrorToJSON(err error) string {
	s, ok := status.FromError(err)
	if !ok {
		return `{"code":0,"message":"not a gRPC error"}`
	}

	errInfo := map[string]interface{}{
		"code":    s.Code(),
		"message": s.Message(),
	}

	jsonBytes, _ := json.Marshal(errInfo)
	return string(jsonBytes)
}
