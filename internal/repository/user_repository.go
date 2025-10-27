package repository

import (
	"ac-ai/internal/models"
	"ac-ai/internal/schemas"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type PaginatedUsersResult struct {
	Users      []models.User
	TotalItems int64
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	// Preload("FinancialProfile") автоматически "джойнит" профиль
	if err := r.db.Preload("FinancialProfile").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("FinancialProfile").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(req *schemas.RegisterRequest, hashedPassword string) (*models.User, error) {
	user := models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	}

	// Используем транзакцию, чтобы создать и юзера, и профиль (если надо)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Создаем пользователя
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 2. Если это CLIENT, создаем профиль
		if user.Role == models.RoleClient {
			if req.ProfileData == nil {
				return fmt.Errorf("profile_data is required for CLIENT role")
			}
			profile := models.FinancialProfile{
				UserID:             user.ID, // Связываем с созданным user.ID
				Income:             req.ProfileData.Income,
				MonthlyPayments:    req.ProfileData.MonthlyPayments,
				CreditHistory:      req.ProfileData.CreditHistory,
				JobExperienceYears: req.ProfileData.JobExperienceYears,
				Age:                req.ProfileData.Age,
				IncomeProof:        req.ProfileData.IncomeProof,
			}
			if err := tx.Create(&profile).Error; err != nil {
				return err
			}
		}
		// Транзакция коммитится автоматически, если нет ошибок
		return nil
	})

	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, err
	}

	// Пере-запрашиваем пользователя с профилем
	return r.GetUserByID(user.ID)
}

func (r *UserRepository) GetUsersByRole(role string, pagination schemas.PaginationQuery) (*PaginatedUsersResult, error) {
	var users []models.User
	var totalItems int64

	baseQuery := r.db.Model(&models.User{}).Where("role = ?", role)

	if err := baseQuery.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	err := baseQuery.
		Preload("FinancialProfile").
		Scopes(PaginateScope(pagination.Page, pagination.Limit)).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return &PaginatedUsersResult{
		Users:      users,
		TotalItems: totalItems,
	}, nil
}
