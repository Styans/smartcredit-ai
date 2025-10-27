package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL                string        `mapstructure:"DATABASE_URL"`
	OpenAIAPIKey               string        `mapstructure:"OPENAI_API_KEY"`
	JWTSecretKey               string        `mapstructure:"JWT_SECRET_KEY"`
	JWTAccessTokenExpireMinutes time.Duration `mapstructure:"JWT_ACCESS_TOKEN_EXPIRE_MINUTES"`
	ServerPort                 string        `mapstructure:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, loading from ENV variables: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	
	// Установка значения по умолчанию для порта
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}
	
	// Установка значения по умолчанию для времени жизни токена
	if cfg.JWTAccessTokenExpireMinutes == 0 {
		cfg.JWTAccessTokenExpireMinutes = 60
	}

	return &cfg, nil
}