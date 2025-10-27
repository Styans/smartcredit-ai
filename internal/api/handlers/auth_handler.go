package handlers

import (
	"ac-ai/internal/auth"
	"ac-ai/internal/models"
	"ac-ai/internal/repository"
	"ac-ai/internal/schemas"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	UserRepo *repository.UserRepository
	JWT      *auth.JWTService
}

func NewAuthHandler(repo *repository.UserRepository, jwt *auth.JWTService) *AuthHandler {
	return &AuthHandler{UserRepo: repo, JWT: jwt}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req schemas.RegisterRequest

	// Валидация JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка на CLIENT и ProfileData
	if req.Role == models.RoleClient && req.ProfileData == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "profile_data is required for CLIENT role"})
		return
	}

	// Проверка, что юзер не существует
	if _, err := h.UserRepo.GetUserByEmail(req.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Хешируем пароль
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Создаем юзера (и профиль)
	user, err := h.UserRepo.CreateUser(&req, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	// (Можно вернуть DTO вместо модели, но для демо сойдет)
	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req schemas.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ищем юзера
	user, err := h.UserRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Проверяем пароль
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Создаем токен
	token, err := h.JWT.CreateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, schemas.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}