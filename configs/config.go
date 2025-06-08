package configs

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Config struct {
	ServerPort int
	JWT
	Database
}

type Database struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

type JWT struct {
	SigningKey string
}

func InitConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("SERVER_PORT", "8000")
	viper.SetDefault("JWT_SIGNING_KEY", "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk")

	cfg := Config{
		ServerPort: viper.GetInt("SERVER_PORT"),
		Database: Database{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			Username: viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		JWT: JWT{
			SigningKey: viper.GetString("JWT_SIGNING_KEY"),
		},
	}
	return &cfg, nil
}
