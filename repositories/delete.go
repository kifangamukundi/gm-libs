package repositories

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

func (r *GormRepository) SoftDelete(model any, id any) error {
	if err := r.DB.Model(model).Where("id = ?", id).Update("deleted_at", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to soft delete model: %v", err)
	}

	return nil
}

func (r *GormRepository) HardDelete(model any, id any) error {
	if err := r.DB.Unscoped().Delete(model, id).Error; err != nil {
		return fmt.Errorf("failed to hard delete model: %v", err)
	}

	return nil
}

func (r *GormRepository) BulkDelete(model any, field string, values []any) error {
	if reflect.ValueOf(model).Kind() != reflect.Pointer {
		return fmt.Errorf("model must be a pointer")
	}
	// if reflect.ValueOf(model).Kind() != reflect.Ptr {
	// 	return fmt.Errorf("model must be a pointer")
	// }

	if field == "" {
		return fmt.Errorf("field cannot be empty")
	}

	if len(values) == 0 {
		return fmt.Errorf("values cannot be empty")
	}

	result := r.DB.Where(field+" IN ?", values).Delete(model)
	if result.Error != nil {
		log.Printf("Error deleting records: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no records found with the provided values")
	}

	return nil
}
