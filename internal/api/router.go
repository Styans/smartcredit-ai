package api

import (
	"ac-ai/internal/api/handlers"
	"ac-ai/internal/api/middleware"
	"ac-ai/internal/auth"
	"ac-ai/internal/config"
	"ac-ai/internal/models"
	"ac-ai/internal/repository"
	"ac-ai/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// ... (Настройка CORS) ...
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// --- ОБНОВЛЕННАЯ ИНИЦИАЛИЗАЦИЯ ---
	// Инициализация сервисов
	userRepo := repository.NewUserRepository(db)
	appRepo := repository.NewApplicationRepository(db) // <-- НОВЫЙ РЕПО
	jwtService := auth.NewJWTService(cfg)
	aiService := services.NewAIService(cfg)

	// Инициализация хэндлеров
	authHandler := handlers.NewAuthHandler(userRepo, jwtService)
	// Передаем appRepo в scoringHandler
	scoringHandler := handlers.NewScoringHandler(userRepo, appRepo, aiService)
	agentHandler := handlers.NewAgentHandler(userRepo, appRepo) // <-- НОВЫЙ ХЭНДЛЕР

	// Группа роутов
	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		scoringGroup := v1.Group("/scoring")
		{
			scoringGroup.Use(middleware.AuthMiddleware(jwtService))
			scoringGroup.Use(middleware.RoleMiddleware(models.RoleClient))
			scoringGroup.POST("/ask", scoringHandler.Ask)
		}

		// --- НОВЫЙ БЛОК: КАБИНЕТ АГЕНТА ---
		agentGroup := v1.Group("/agent")
		{
			agentGroup.Use(middleware.AuthMiddleware(jwtService))
			agentGroup.Use(middleware.RoleMiddleware(models.RoleAgent))

			// Дашборд: Заявки на ручное рассмотрение
			agentGroup.GET("/applications/review", agentHandler.GetApplicationsForReview)
			// Мониторинг: Все клиенты
			agentGroup.GET("/clients", agentHandler.GetAllClients)
			// Мониторинг: Все заявки
			agentGroup.GET("/applications/all", agentHandler.GetAllApplications)
			// Мониторинг: Все клиенты
		}
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to AC-AI Scoring API (Go Version)"})
	})

	return r
}
