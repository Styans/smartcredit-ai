package schemas

import (
	"time"
)

// Упрощенная информация о пользователе для списка заявок
type ApplicationUserOut struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

// Полная информация о заявке для агента
type ApplicationOut struct {
	ID              uint               `json:"id"`
	CreatedAt       time.Time          `json:"created_at"`
	User            ApplicationUserOut `json:"user"` // Вложенный пользователь
	RequestedAmount float64            `json:"requested_amount"`
	FinalDecision   string             `json:"final_decision"` // Решение ИИ
	ColdScore       int                `json:"cold_score"`
	AIResponse      string             `json:"ai_response"`  // Что увидел клиент
	AgentStatus     string             `json:"agent_status"` // Статус от агента
	AgentNotes      string             `json:"agent_notes"`
	InternalReasons []string           `json:"internal_reasons"`
}

// Профиль клиента для просмотра агентом
type ClientProfileOut struct {
	ID               uint                   `json:"id"`
	Email            string                 `json:"email"`
	FinancialProfile FinancialProfileCreate `json:"financial_profile"`
	// (можно добавить историю заявок этого клиента)
}

// PaginationQuery - параметры, которые мы ожидаем из URL (?page=1&limit=10)
type PaginationQuery struct {
	Page  int `form:"page,default=1"`   // 'form' тэг для Gin
	Limit int `form:"limit,default=10"` // 'form' тэг для Gin
}

// PaginationMeta - информация о пагинации для фронтенда
type PaginationMeta struct {
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	CurrentPage  int   `json:"current_page"`
	ItemsPerPage int   `json:"items_per_page"`
}

// PaginatedResponse - общий контейнер для ответа с пагинацией
// Мы используем дженерики (any), чтобы переиспользовать эту схему
type PaginatedResponse struct {
	Data any            `json:"data"` // Здесь будут лежать наши заявки или клиенты
	Meta PaginationMeta `json:"meta"`
}
