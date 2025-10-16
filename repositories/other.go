package repositories

import (
	"fmt"

	"gorm.io/gorm"
)

func (r *GormRepository) Count(query *gorm.DB, model any) (int64, error) {
	var count int64

	if err := query.Model(model).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error counting models: %v", err)
	}

	return count, nil
}

func (r *GormRepository) CountConditional(model any, conditions map[string]any) (int64, error) {
	var count int64

	query := r.DB.Model(model)

	if len(conditions) > 0 {
		query = query.Where(conditions)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error counting models: %v", err)
	}

	return count, nil
}

func (r *GormRepository) Find(query *gorm.DB, model any) ([]any, error) {
	var result []any

	if err := query.Find(&result).Error; err != nil {
		return nil, fmt.Errorf("error fetching models: %v", err)
	}

	return result, nil
}
