package repository

import (
	"ac-ai/internal/models"
	"ac-ai/internal/schemas"

	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

type PaginatedApplicationsResult struct {
	Applications []models.ScoringApplication
	TotalItems   int64
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

// CreateApplication - Вызывается хэндлером клиента при подаче заявки
func (r *ApplicationRepository) CreateApplication(app *models.ScoringApplication) error {
	return r.db.Create(app).Error
}

// GetApplicationsForReview - Вызывается агентом (главный дашборд)
// Показывает заявки, требующие ручного решения
func (r *ApplicationRepository) GetApplicationsForReview(pagination schemas.PaginationQuery) (*PaginatedApplicationsResult, error) {
	var applications []models.ScoringApplication
	var totalItems int64

	// Сначала считаем общее количество (для пагинации)
	baseQuery := r.db.Model(&models.ScoringApplication{}).
		Where("final_decision = ? AND agent_status = ?", models.StatusManualReview, models.AgentStatusPending)

	if err := baseQuery.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	// Теперь получаем нужную "страницу"
	err := baseQuery.
		Preload("User").
		Order("created_at desc").
		Scopes(PaginateScope(pagination.Page, pagination.Limit)). // <-- Применяем пагинацию
		Find(&applications).Error

	if err != nil {
		return nil, err
	}

	return &PaginatedApplicationsResult{
		Applications: applications,
		TotalItems:   totalItems,
	}, nil
}

// GetAllApplications - Вызывается агентом (общий мониторинг)
func (r *ApplicationRepository) GetAllApplications(pagination schemas.PaginationQuery) (*PaginatedApplicationsResult, error) {
	var applications []models.ScoringApplication
	var totalItems int64

	baseQuery := r.db.Model(&models.ScoringApplication{})
	
	if err := baseQuery.Count(&totalItems).Error; err != nil {
		return nil, err
	}
	
	err := baseQuery.
		Preload("User").
		Order("created_at desc").
		Scopes(PaginateScope(pagination.Page, pagination.Limit)).
		Find(&applications).Error
		
	if err != nil {
		return nil, err
	}

	return &PaginatedApplicationsResult{
		Applications: applications,
		TotalItems:   totalItems,
	}, nil
}