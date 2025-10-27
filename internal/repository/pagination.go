// internal/repository/pagination.go
package repository

import (
	"math"
	"gorm.io/gorm"
)

// PaginateScope - это GORM Scope, который можно переиспользовать
// Он принимает страницу и лимит и применяет Offset/Limit к запросу
func PaginateScope(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		} else if limit > 100 {
			limit = 100 // Защита от слишком больших запросов
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

// CalculateMeta - вычисляет мета-данные пагинации
func CalculateMeta(totalItems int64, page, limit int) (int, int) {
	if limit <= 0 {
		limit = 10
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if page <= 0 {
		page = 1
	}
	return totalPages, page
}