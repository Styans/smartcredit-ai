package models

import (
	"gorm.io/gorm"
)

// Роли будем хранить как строки для простоты
const (
	RoleClient = "CLIENT"
	RoleAgent  = "AGENT"
)

type User struct {
	gorm.Model
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"type:varchar(10);not null"`

	FinancialProfile FinancialProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	// --- ДОБАВЬТЕ ЭТУ СТРОКУ ---
	ScoringApplications []ScoringApplication `gorm:"foreignKey:UserID"`
}
