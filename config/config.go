package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl      string
	ServerPort string
	Region     string
	BucketName string
	AccessId   string
	AcessKey   string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	config := Config{
		DBUrl:      os.Getenv("DB_URL"),
		ServerPort: os.Getenv("SERVER_PORT"),
		Region:     os.Getenv("REGION"),
		BucketName: os.Getenv("BUCKET"),
		AccessId:   os.Getenv("AWS_ACCESS_KEY_ID"),
		AcessKey:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	return &config, nil
}
