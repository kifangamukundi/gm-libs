package repositories

import (
	"fmt"
)

func (r *GormRepository) ClearAssociations(model any, association string) error {
	if err := r.DB.Model(model).Association(association).Clear(); err != nil {
		return fmt.Errorf("failed to clear association %s: %v", association, err)
	}

	return nil
}

func (r *GormRepository) CreateAssociations(model any, association string, associations []any) error {
	tx := r.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// First clear existing associations
	if err := tx.Model(model).Association(association).Clear(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear associations: %v", err)
	}

	// For join tables with extra fields, we need to create them directly
	for _, assoc := range associations {
		// Try to append first (works for normal associations)
		if err := tx.Model(model).Association(association).Append(assoc); err != nil {
			// If append fails, try to create the join table record directly
			if err := tx.Create(assoc).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create association: %v", err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
