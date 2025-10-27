package schemas

// В Go мы используем struct tags для валидации JSON

type FinancialProfileCreate struct {
	Income             float64 `json:"income" binding:"required,gte=0"`
	MonthlyPayments    float64 `json:"monthly_payments" binding:"required,gte=0"`
	CreditHistory      string  `json:"credit_history" binding:"required,oneof=no_issues minor_issues major_issues"`
	JobExperienceYears float64 `json:"job_experience_years" binding:"required,gte=0"`
	Age                int     `json:"age" binding:"required,gte=18"`
	IncomeProof        string  `json:"income_proof" binding:"required,oneof=official indirect verbal"`
}

type RegisterRequest struct {
	Email       string                  `json:"email" binding:"required,email"`
	Password    string                  `json:"password" binding:"required,min=6"`
	Role        string                  `json:"role" binding:"required,oneof=CLIENT AGENT"`
	ProfileData *FinancialProfileCreate `json:"profile_data,omitempty"` // omitempty, т.к. для AGENT его нет
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// ... (Добавьте UserOut, ProfileOut, если нужно)