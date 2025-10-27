package handlers

import (
	"ac-ai/internal/models"
	"ac-ai/internal/repository"
	"ac-ai/internal/schemas"
	"ac-ai/internal/services"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScoringHandler struct {
	AppRepo   *repository.ApplicationRepository
	UserRepo  *repository.UserRepository
	AIService *services.AIService
}

func NewScoringHandler(
	appRepo *repository.ApplicationRepository,
	repo *repository.UserRepository,
	ai *services.AIService,
) *ScoringHandler {
	return &ScoringHandler{
		AppRepo:   appRepo,
		UserRepo:  repo,
		AIService: ai,
	}
}

func (h *ScoringHandler) Ask(c *gin.Context) {
	var req schemas.ScoringRequest

	// 1. Валидация запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Получаем ID юзера из middleware
	userID, _ := c.Get("userID")

	// 3. Получаем полный профиль юзера из БД
	user, err := h.UserRepo.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}
	if user.Role != models.RoleClient || user.FinancialProfile.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not a client or has no profile"})
		return
	}

	ctx := context.Background()

	// 4. Парсим сумму из запроса
	requestedAmount, err := h.AIService.ParseAmountFromQuery(ctx, req.Query)
	if err != nil || requestedAmount == 0 {
		answer := "Я могу помочь с расчетом кредита. Пожалуйста, укажите желаемую сумму, например: 'Хочу 15 000 000 тенге'."
		c.JSON(http.StatusOK, schemas.ScoringResponse{Answer: answer})
		return
	}

	// 5. "Холодный" скоринг
	scoreResult := services.CalculateColdScore(&user.FinancialProfile, requestedAmount)

	// 6. "Теплый" AI-анализ
	answer, err := h.AIService.GetAIAnalysis(ctx, scoreResult, &user.FinancialProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI analysis", "details": err.Error()})
		return
	}

	internalReasonsBytes, _ := json.Marshal(scoreResult.Recommendations)
	internalReasonsStr := string(internalReasonsBytes)

	application := models.ScoringApplication{
		UserID:          user.ID,
		RequestedAmount: requestedAmount,
		FinalDecision:   scoreResult.Decision,
		ColdScore:       scoreResult.TotalScore,
		AIResponse:      answer,
		InternalReasons: internalReasonsStr,
		AgentStatus:     models.AgentStatusPending, // По умолчанию ждет
	}

	// Если решение НЕ ручное, то агенту не нужно ничего делать
	if application.FinalDecision != models.StatusManualReview {
		application.AgentStatus = application.FinalDecision
	}

	// 8. Сохраняем в БД
	if err := h.AppRepo.CreateApplication(&application); err != nil {
		// Не показываем ошибку клиенту, но логируем ее
		log.Printf("CRITICAL: Failed to save application for user %d: %v", user.ID, err)
	}

	// --- ** КОНЕЦ НОВОЙ ЛОГИКИ ** ---

	// 9. Отправляем ответ клиенту
	c.JSON(http.StatusOK, schemas.ScoringResponse{Answer: answer})
}
