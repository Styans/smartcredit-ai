// internal/models/application.go
package models

import (
	"gorm.io/gorm"
)

// Статусы
const (
	StatusApproved     = "APPROVED"
	StatusDenied       = "DENIED"
	StatusManualReview = "MANUAL_REVIEW"
)

// Статусы Агента
const (
	AgentStatusPending  = "PENDING"
	AgentStatusApproved = "AGENT_APPROVED"
	AgentStatusDenied   = "AGENT_DENIED"
)

type ScoringApplication struct {
	gorm.Model
	UserID          uint    `gorm:"not null"`
	RequestedAmount float64 `gorm:"not null"`

	// Решение, которое принял ИИ / "холодный" скоринг
	FinalDecision   string `gorm:"type:varchar(20);not null"` // APPROVED, DENIED, MANUAL_REVIEW
	ColdScore       int
	AIResponse      string `gorm:"type:text"` // Ответ, который увидел клиент
	// Поля для Агента
	InternalReasons string `gorm:"type:text"` // JSON-массив []string с причинами для агента
	
	AgentStatus string `gorm:"type:varchar(20);default:'PENDING'"` // Статус, который выставил агент
	AgentNotes  string `gorm:"type:text"`                          // Комментарий агента

	User User `gorm:"foreignKey:UserID"` // Связь с пользователем
}
