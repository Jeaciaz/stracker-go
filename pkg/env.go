package pkg

import (
	"os"
	_ "github.com/joho/godotenv/autoload"
)

type Env struct {
	DbUrl    string
	Password string
}

func LoadConfig() (*Env, error) {
	cfg := &Env{
		DbUrl: os.Getenv("DB_URL"),
		Password: os.Getenv("PASSWORD"),
	}

	return cfg, nil
}
