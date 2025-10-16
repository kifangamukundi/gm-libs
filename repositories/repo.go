package repositories

import (
	"gorm.io/gorm"
)

type Repository interface {
	Begin() (Repository, error)
	Commit() error
	Rollback() error
	WithTransaction(tx *gorm.DB) Repository

	Create(model any) error

	GetByID(model any, id any) (any, error)
	GetAll(model any) (any, error)
	GetByField(model any, field string, value any) (any, error)
	GetByFields(model any, fieldValues map[string]any) (any, error)
	GetByFieldWithPreload(model any, field, value string, preload ...string) (any, error)
	GetAllByFieldWithPreload(model any, field, value string, preload ...string) (any, error)
	GetAllWithPreload(model any, preload ...string) (any, error)
	GetMonthCount(model any) ([]map[string]any, error)

	Update(model any) error

	SoftDelete(model any, id any) error
	HardDelete(model any, id any) error
	BulkDelete(model any, field string, values []any) error

	ClearAssociations(model any, association string) error
	CreateAssociations(model any, association string, associations []any) error

	GetFilteredByItems(field string, items []any, model any, preload []string) ([]any, error)
	GetAllFiltered(skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria any, model any, preload []string) ([]any, int64, int64, error)
	GetAllFilteredByField(
		model any,
		filters map[string]any,
		preload []string,
		joins []struct {
			Table     string
			Condition string
			JoinType  string
		},
	) ([]any, error)

	Count(query *gorm.DB, model any) (int64, error)
	CountConditional(model any, conditions map[string]any) (int64, error)
	Find(query *gorm.DB, model any) ([]any, error)
}

type GormRepository struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{DB: db}
}
