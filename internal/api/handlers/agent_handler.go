package handlers

import (
	"ac-ai/internal/models"
	"ac-ai/internal/repository"
	"ac-ai/internal/schemas"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	UserRepo *repository.UserRepository
	AppRepo  *repository.ApplicationRepository
}

func NewAgentHandler(userRepo *repository.UserRepository, appRepo *repository.ApplicationRepository) *AgentHandler {
	return &AgentHandler{
		UserRepo: userRepo,
		AppRepo:  appRepo,
	}
}

// GET /api/v1/agent/applications/review
func (h *AgentHandler) GetApplicationsForReview(c *gin.Context) {
	// 1. Парсим параметры пагинации из URL (?page=1&limit=10)
	var pagination schemas.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	// 2. Получаем данные из репозитория
	result, err := h.AppRepo.GetApplicationsForReview(pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	// 3. Конвертируем в DTO (Data Transfer Object)
	var applicationsOut []schemas.ApplicationOut
	for _, app := range result.Applications {
		// Десериализуем InternalReasons из JSON-строки в []string
		var reasons []string
		if app.InternalReasons != "" {
			// Игнорируем ошибку, если JSON невалидный, просто вернется пустой слайс
			_ = json.Unmarshal([]byte(app.InternalReasons), &reasons)
		}

		applicationsOut = append(applicationsOut, schemas.ApplicationOut{
			ID:              app.ID,
			CreatedAt:       app.CreatedAt,
			User:            schemas.ApplicationUserOut{ID: app.User.ID, Email: app.User.Email},
			RequestedAmount: app.RequestedAmount,
			FinalDecision:   app.FinalDecision,
			ColdScore:       app.ColdScore,
			AIResponse:      app.AIResponse,
			AgentStatus:     app.AgentStatus,
			AgentNotes:      app.AgentNotes,
			InternalReasons: reasons, // <-- Передаем []string
		})
	}

	// 4. Считаем мета-данные пагинации
	totalPages, currentPage := repository.CalculateMeta(result.TotalItems, pagination.Page, pagination.Limit)

	// 5. Возвращаем обернутый ответ
	c.JSON(http.StatusOK, schemas.PaginatedResponse{
		Data: applicationsOut,
		Meta: schemas.PaginationMeta{
			TotalItems:   result.TotalItems,
			TotalPages:   totalPages,
			CurrentPage:  currentPage,
			ItemsPerPage: pagination.Limit,
		},
	})
}

// GET /api/v1/agent/clients - МОДИФИЦИРОВАНО
func (h *AgentHandler) GetAllClients(c *gin.Context) {
	var pagination schemas.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	result, err := h.UserRepo.GetUsersByRole(models.RoleClient, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clients"})
		return
	}

	var clientsOut []schemas.ClientProfileOut
	for _, user := range result.Users {
		if user.FinancialProfile.ID == 0 {
			continue
		}
		clientsOut = append(clientsOut, schemas.ClientProfileOut{
			ID:    user.ID,
			Email: user.Email,
			FinancialProfile: schemas.FinancialProfileCreate{
				Income:             user.FinancialProfile.Income,
				MonthlyPayments:    user.FinancialProfile.MonthlyPayments,
				CreditHistory:      user.FinancialProfile.CreditHistory,
				JobExperienceYears: user.FinancialProfile.JobExperienceYears,
				Age:                user.FinancialProfile.Age,
				IncomeProof:        user.FinancialProfile.IncomeProof,
			},
		})
	}

	totalPages, currentPage := repository.CalculateMeta(result.TotalItems, pagination.Page, pagination.Limit)

	c.JSON(http.StatusOK, schemas.PaginatedResponse{
		Data: clientsOut,
		Meta: schemas.PaginationMeta{
			TotalItems:   result.TotalItems,
			TotalPages:   totalPages,
			CurrentPage:  currentPage,
			ItemsPerPage: pagination.Limit,
		},
	})
}

func (h *AgentHandler) GetAllApplications(c *gin.Context) {
	var pagination schemas.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	result, err := h.AppRepo.GetAllApplications(pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all applications"})
		return
	}

	var applicationsOut []schemas.ApplicationOut
	for _, app := range result.Applications {
		var reasons []string
		if app.InternalReasons != "" {
			_ = json.Unmarshal([]byte(app.InternalReasons), &reasons)
		}
		applicationsOut = append(applicationsOut, schemas.ApplicationOut{
			ID:              app.ID,
			CreatedAt:       app.CreatedAt,
			User:            schemas.ApplicationUserOut{ID: app.User.ID, Email: app.User.Email},
			RequestedAmount: app.RequestedAmount,
			FinalDecision:   app.FinalDecision,
			ColdScore:       app.ColdScore,
			AIResponse:      app.AIResponse,
			AgentStatus:     app.AgentStatus,
			AgentNotes:      app.AgentNotes,
			InternalReasons: reasons,
		})
	}

	totalPages, currentPage := repository.CalculateMeta(result.TotalItems, pagination.Page, pagination.Limit)

	c.JSON(http.StatusOK, schemas.PaginatedResponse{
		Data: applicationsOut,
		Meta: schemas.PaginationMeta{
			TotalItems:   result.TotalItems,
			TotalPages:   totalPages,
			CurrentPage:  currentPage,
			ItemsPerPage: pagination.Limit,
		},
	})
}