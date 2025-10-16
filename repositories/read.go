package repositories

import (
	"fmt"
	"log"
	"reflect"
)

func (r *GormRepository) GetByID(model any, id any) (any, error) {
	if err := r.DB.Where("id = ?", id).First(model).Error; err != nil {
		return nil, fmt.Errorf("failed to get model by ID: %v", err)
	}

	return model, nil
}

func (r *GormRepository) GetAll(model any) (any, error) {
	if err := r.DB.Find(model).Error; err != nil {
		return nil, fmt.Errorf("failed to get all items: %v", err)
	}

	return model, nil
}

func (r *GormRepository) GetByField(model any, field string, value any) (any, error) {
	if err := r.DB.Where(fmt.Sprintf("%s = ?", field), value).First(model).Error; err != nil {
		return nil, fmt.Errorf("failed to get model by %s: %v", field, err)
	}

	return model, nil
}

func (r *GormRepository) GetByFields(model any, fieldValues map[string]any) (any, error) {
	if len(fieldValues) == 0 {
		return nil, fmt.Errorf("field values cannot be empty")
	}

	query := r.DB
	for field, value := range fieldValues {
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			query = query.Where(fmt.Sprintf("%s IN (?)", field), value)
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	if err := query.Find(model).Error; err != nil {
		return nil, fmt.Errorf("failed to get model by fields: %v", err)
	}

	return model, nil
}

func (r *GormRepository) GetByFieldWithPreload(model any, field, value string, preload ...string) (any, error) {
	query := r.DB.Where(field+" = ?", value)

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.First(model).Error
	if err != nil {
		log.Printf("Error fetching record by %s: %v", field, err)
		return nil, err
	}

	return model, nil
}

func (r *GormRepository) GetAllByFieldWithPreload(model any, field, value string, preload ...string) (any, error) {
	query := r.DB.Where(field+" = ?", value)

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Find(model).Error
	if err != nil {
		log.Printf("Error fetching record by %s: %v", field, err)
		return nil, err
	}

	return model, nil
}

func (r *GormRepository) GetAllWithPreload(model any, preload ...string) (any, error) {
	query := r.DB

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Find(model).Error
	if err != nil {
		log.Printf("Error fetching records: %v", err)
		return nil, err
	}

	return model, nil
}

func (r *GormRepository) GetMonthCount(model any) ([]map[string]any, error) {
	var results []map[string]any

	err := r.DB.Model(model).
		Select("TO_CHAR(created_at, 'Mon') AS month, COUNT(*) AS count, EXTRACT(MONTH FROM created_at) AS month_num").
		Group("TO_CHAR(created_at, 'Mon'), month_num").
		Order("month_num").
		Scan(&results).Error
	if err != nil {
		log.Printf("Error fetching monthly counts: %v", err)
		return nil, err
	}

	return results, nil
}
