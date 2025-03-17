package utils

import (
	"2025_CakeLand_API/internal/pkg/config"
	"database/sql"
	"fmt"
)

func ConnectPostgres(cfg *config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	pingErr := db.Ping()
	return db, pingErr
}
