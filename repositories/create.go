package repositories

import "fmt"

func (r *GormRepository) Create(model any) error {
	if err := r.DB.Create(model).Error; err != nil {
		return fmt.Errorf("failed to create model: %v", err)
	}

	return nil
}
