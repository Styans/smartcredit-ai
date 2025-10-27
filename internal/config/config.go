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
	// --- ИСПРАВЛЕНИЕ ЗДЕСЬ ---
	// 1. Сначала говорим Viper явно искать эти ENV-переменные
	// Это гарантирует, что переменные Render будут найдены
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("OPENAI_API_KEY")
	viper.BindEnv("JWT_SECRET_KEY")
	viper.BindEnv("JWT_ACCESS_TOKEN_EXPIRE_MINUTES")
	viper.BindEnv("SERVER_PORT")

	// 2. Затем (для локальной разработки) пытаемся прочитать .env
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, loading *only* from ENV variables.")
	}
    // --- КОНЕЦ ИСПРАВЛЕНИЯ ---

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Установка значения по умолчанию для порта
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080" // Render ожидает порт 8080 или 10000
	}

	// Установка значения по умолчанию для времени жизни токена
	if cfg.JWTAccessTokenExpireMinutes == 0 {
		// Важно: Viper не парсит '60m' из .env, если он не нашел BindEnv
		// Поэтому мы устанавливаем '60' и умножаем
		cfg.JWTAccessTokenExpireMinutes = 60 * time.Minute
	}

	return &cfg, nil
}