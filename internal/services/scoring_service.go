package services

import (
	"ac-ai/internal/models"
)

// ** НОВОЕ ПОЛЕ **
// Мы добавили DtiRatio и RecommendedMaxAmount
type ColdScoreResult struct {
	TotalScore           int
	Decision             string // "APPROVED", "DENIED", "MANUAL_REVIEW"
	DtiRatio             float64
	RecommendedMaxAmount float64 // Максимальная сумма, которую мы можем рекомендовать
	RequestedAmount      float64
	Recommendations      []string
}

// Константа для "идеальной" долговой нагрузки (например, 40%)
const maxSafeDTI = 0.40

// Срок кредита по умолчанию для расчета (60 мес = 5 лет)
const defaultLoanTermMonths = 60

func CalculateColdScore(profile *models.FinancialProfile, requestedAmount float64) *ColdScoreResult {
	baseScore := 0
	recommendations := []string{}

	// --- ** НОВАЯ ЛОГИКА ** ---
	// Рассчитываем, сколько пользователь может платить в месяц
	maxTotalMonthlyPayment := profile.Income * maxSafeDTI
	availableForNewPayment := maxTotalMonthlyPayment - profile.MonthlyPayments

	// Если он уже тратит слишком много, он не может позволить себе новый кредит
	if availableForNewPayment < 0 {
		availableForNewPayment = 0
	}

	// Рассчитываем максимальную сумму, которую он может взять
	// (Это обратный расчет от ежемесячного платежа)
	recommendedMaxAmount := availableForNewPayment * defaultLoanTermMonths

	// --- Конец новой логики ---

	// 1. DTI (Долговая нагрузка)
	var dti float64
	newMonthlyPayment := requestedAmount / defaultLoanTermMonths
	totalPayments := profile.MonthlyPayments + newMonthlyPayment

	if profile.Income > 0 {
		dti = totalPayments / profile.Income
	} else {
		dti = 1.0 // Плохой DTI, если доход 0
	}

	if dti < 0.2 {
		baseScore += 300
	} else if dti < 0.4 {
		baseScore += 150
	} else if dti < 0.6 {
		baseScore += 50
	} else {
		baseScore -= 100
		// Добавляем конкретную причину
		recommendations = append(recommendations, "Долговая нагрузка (DTI) слишком высока.")
	}

	// 2. Кредитная история
	switch profile.CreditHistory {
	case models.CreditHistoryNoIssues:
		baseScore += 300
	case models.CreditHistoryMinorIssues:
		baseScore += 100
	case models.CreditHistoryMajorIssues:
		baseScore -= 200
		recommendations = append(recommendations, "Плохая кредитная история является негативным фактором.")
	}

	// 3. Стаж работы
	if profile.JobExperienceYears > 3 {
		baseScore += 200
	} else if profile.JobExperienceYears >= 1 {
		baseScore += 100
	} else {
		recommendations = append(recommendations, "Стаж работы менее 1 года - это фактор риска.")
	}

	// 4. Проверка на СВЕРХ-сумму
	// Если запрошенная сумма (50 млрд) В РАЗЫ больше, чем мы можем дать,
	// это гарантированный отказ или ручная проверка.
	if requestedAmount > recommendedMaxAmount*1.5 { // Если просят на 50% больше, чем можно
		// Если DTI был в порядке (например, 20 млн / 60 = высокий платеж),
		// но РЕКОМЕНДУЕМАЯ сумма (например, 5 млн) намного ниже,
		// мы принудительно снижаем балл и меняем решение.
		baseScore -= 200 // Сильный штраф
		recommendations = append(recommendations, "Запрошенная сумма значительно превышает ваши финансовые возможности.")
	}

	var decision string
	if baseScore < 400 {
		decision = "DENIED"
	} else if baseScore < 700 {
		decision = "MANUAL_REVIEW"
	} else {
		decision = "APPROVED"
	}

	// Если финальное решение DENIED или MANUAL, а DTI был плохой,
	// мы принудительно ставим решение DENIED.
	if (decision == "DENIED" || decision == "MANUAL_REVIEW") && dti > 0.6 {
		decision = "DENIED"
	}

	return &ColdScoreResult{
		TotalScore:           baseScore,
		Decision:             decision,
		DtiRatio:             dti,                  // ** Добавили **
		RecommendedMaxAmount: recommendedMaxAmount, // ** Добавили **
		RequestedAmount:      requestedAmount,
		Recommendations:      recommendations,
	}
}
