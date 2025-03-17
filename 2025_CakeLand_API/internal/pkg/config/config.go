package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env   string         `yaml:"env" env-default:"local"`
	GRPC  GRPCConfig     `yaml:"grpc"`
	DB    DatabaseConfig `yaml:"database"`
	MinIO MinioConfig    `yaml:"minio"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DatabaseConfig struct {
	Host     string
	DBName   string
	Port     int
	User     string
	Password string
	SSLMode  string
}

type MinioConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Host      string `json:"host"`
	Bucket    string `json:"bucket"`
	Region    string `json:"region"`
	UseSSL    bool   `json:"use_ssl"`
}

func NewConfig() (*Config, error) {
	configPath := fetchConfigPath()
	if configPath == "" {
		return nil, fmt.Errorf("путь к конфигурационному файлу пустой")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл конфигурации не существует: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %s", err.Error())
	}

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Ошибка при загрузке .env файла:", err)
	}

	// Чтение переменных окружения для БД
	cfg.DB.Host = os.Getenv("HOST")
	cfg.DB.DBName = os.Getenv("DB_NAME")
	cfg.DB.Port, _ = strconv.Atoi(os.Getenv("PORT"))
	cfg.DB.User = os.Getenv("POSTGRES_USER")
	cfg.DB.Password = os.Getenv("PASSWORD")
	cfg.DB.SSLMode = os.Getenv("SSL_MODE")

	// Чтение переменных окружения для MinIO
	cfg.MinIO.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	cfg.MinIO.SecretKey = os.Getenv("MINIO_SECRET_KEY")
	cfg.MinIO.Host = os.Getenv("MINIO_HOST")
	cfg.MinIO.Bucket = os.Getenv("MINIO_BUCKET")
	cfg.MinIO.Region = os.Getenv("MINIO_REGION")
	useSSL := os.Getenv("MINIO_USE_SSL")
	if useSSL == "true" {
		cfg.MinIO.UseSSL = true
	} else {
		cfg.MinIO.UseSSL = false
	}

	return &cfg, nil
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "путь к конфигурационному файлу")
	flag.Parse()
	if res == "" {
		res = "./config/config.yaml"
	}
	return res
}
