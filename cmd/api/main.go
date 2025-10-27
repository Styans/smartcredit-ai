package main

import (
	"ac-ai/internal/api"
	"ac-ai/internal/config"
	"ac-ai/internal/database"
	"log"
)

func main() {
	// 1. Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	
	// 2. Подключение к БД
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// 3. Настройка роутера
	router := api.SetupRouter(db, cfg)

	// 4. Запуск сервера
	log.Printf("Starting server on port %s...", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}