package repositories

import "fmt"

func (r *GormRepository) Update(model any) error {
	if err := r.DB.Save(model).Error; err != nil {
		return fmt.Errorf("failed to update model: %v", err)
	}

	return nil
}
