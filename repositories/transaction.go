package repositories

import (
	"gorm.io/gorm"
)

func (r *GormRepository) Begin() (Repository, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &GormRepository{DB: tx}, nil
}

func (r *GormRepository) Commit() error {
	return r.DB.Commit().Error
}

func (r *GormRepository) Rollback() error {
	return r.DB.Rollback().Error
}

func (r *GormRepository) WithTransaction(tx *gorm.DB) Repository {
	return &GormRepository{DB: tx}
}
