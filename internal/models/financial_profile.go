package models

import (
	"gorm.io/gorm"
)

// Типы будем хранить как строки
const (
	CreditHistoryNoIssues    = "no_issues"
	CreditHistoryMinorIssues = "minor_issues"
	CreditHistoryMajorIssues = "major_issues"

	IncomeProofOfficial = "official"
	IncomeProofIndirect = "indirect"
	IncomeProofVerbal   = "verbal"
)

type FinancialProfile struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex;not null"`

	Income               float64 `gorm:"not null"`
	MonthlyPayments      float64 `gorm:"not null"`
	CreditHistory        string  `gorm:"type:varchar(20);not null"`
	JobExperienceYears float64 `gorm:"not null"`
	Age                  int     `gorm:"not null"`
	IncomeProof          string  `gorm:"type:varchar(20);not null"`
}