package repositories

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kifangamukundi/gm-libs/queryparams"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (r *GormRepository) GetFilteredByItems(field string, items []any, model any, preload []string) ([]any, error) {
	var result []any

	query := r.DB.Model(model)

	for _, assoc := range preload {
		query = query.Preload(assoc)
	}

	if len(items) > 0 {
		query = query.Where(fmt.Sprintf("%s IN (?)", field), items)
	}

	resultPtr := reflect.New(reflect.SliceOf(reflect.TypeOf(model)))
	if err := query.Find(resultPtr.Interface()).Error; err != nil {
		return nil, fmt.Errorf("error fetching models: %v", err)
	}

	resultSlice := resultPtr.Elem()
	for i := range resultSlice.Len() {
		result = append(result, resultSlice.Index(i).Interface())
	}

	return result, nil
}

func (r *GormRepository) GetAllFiltered(skip, limit int, sortOrder, sortByColumn, searchRegex string, searchColumns []string, filterCriteria any, model any, preload []string) ([]any, int64, int64, error) {
	var result []any

	query := r.DB.Model(model)

	for _, assoc := range preload {
		query = query.Preload(assoc)
	}

	if searchRegex != "" {
		searchCondition := queryparams.BuildSearchCondition(searchColumns, searchRegex)
		query = query.Where(searchCondition)
	}

	if filterCriteria != nil {
		filters := filterCriteria.(map[string]any)

		for key, value := range filters {
			switch key {
			case "created_at_gte":
				mainTableName := r.DB.NamingStrategy.TableName(reflect.TypeOf(model).Elem().Name())
				query = query.Where(fmt.Sprintf("%s.created_at >= ?", mainTableName), value)
			case "created_at_lte":
				mainTableName := r.DB.NamingStrategy.TableName(reflect.TypeOf(model).Elem().Name())
				query = query.Where(fmt.Sprintf("%s.created_at <= ?", mainTableName), value)
			case "updated_at_gte":
				mainTableName := r.DB.NamingStrategy.TableName(reflect.TypeOf(model).Elem().Name())
				query = query.Where(fmt.Sprintf("%s.updated_at >= ?", mainTableName), value)
			case "updated_at_lte":
				mainTableName := r.DB.NamingStrategy.TableName(reflect.TypeOf(model).Elem().Name())
				query = query.Where(fmt.Sprintf("%s.updated_at <= ?", mainTableName), value)
			default:
				if relatedFilters, ok := value.(map[string]any); ok {
					modelType := reflect.TypeOf(model)
					field, ok := modelType.Elem().FieldByName(cases.Title(language.English).String(key))
					if !ok {
						return nil, 0, 0, fmt.Errorf("relationship '%s' not found", key)
					}

					relationshipType := field.Type
					if relationshipType.Kind() == reflect.Slice {
						relationshipType = relationshipType.Elem()
					}

					relatedTableName := r.DB.NamingStrategy.TableName(relationshipType.Name())
					mainTableName := r.DB.NamingStrategy.TableName(modelType.Elem().Name())

					var joinClause string
					foreignKeyTag := field.Tag.Get("gorm")

					if strings.Contains(foreignKeyTag, "many2many:") {
						pivotTableName := strings.Split(foreignKeyTag, ":")[1]
						joinClause = fmt.Sprintf(
							"JOIN %s ON %s.%s_id = %s.id JOIN %s ON %s.id = %s.%s_id",
							pivotTableName, pivotTableName, strings.ToLower(modelType.Elem().Name()), mainTableName,
							relatedTableName, relatedTableName, pivotTableName, strings.ToLower(relationshipType.Name()),
						)
					} else if strings.Contains(foreignKeyTag, "foreignKey:") {
						foreignKeyColumn := r.DB.NamingStrategy.ColumnName(mainTableName, strings.Split(foreignKeyTag, ":")[1])
						joinClause = fmt.Sprintf("JOIN %s ON %s.%s = %s.id", relatedTableName, mainTableName, foreignKeyColumn, relatedTableName)
					} else {
						foreignKeyColumn := r.DB.NamingStrategy.ColumnName(mainTableName, fmt.Sprintf("%s_id", strings.ToLower(relatedTableName)))
						joinClause = fmt.Sprintf("JOIN %s ON %s.%s = %s.id", relatedTableName, mainTableName, foreignKeyColumn, relatedTableName)
					}

					query = query.Joins(joinClause)

					for relatedKey, relatedValue := range relatedFilters {
						query = query.Where(fmt.Sprintf("%s.%s = ?", relatedTableName, relatedKey), relatedValue)
					}
				} else {
					if boolValue, ok := value.(bool); ok {
						query = query.Where(fmt.Sprintf("%s = ?", key), boolValue)
					} else {
						query = query.Where(fmt.Sprintf("%s = ?", key), value)
					}
				}
			}
		}
	}

	totalCount, err := r.Count(query, model)
	if err != nil {
		return nil, 0, 0, err
	}

	filteredCount, err := r.Count(query, model)
	if err != nil {
		return nil, 0, 0, err
	}

	query = query.Order(fmt.Sprintf("%s %s", sortByColumn, sortOrder)).Offset(skip).Limit(limit)

	resultPtr := reflect.New(reflect.SliceOf(reflect.TypeOf(model)))
	if err := query.Find(resultPtr.Interface()).Error; err != nil {
		return nil, 0, 0, fmt.Errorf("error fetching models: %v", err)
	}

	resultSlice := resultPtr.Elem()
	for i := range resultSlice.Len() {
		result = append(result, resultSlice.Index(i).Interface())
	}

	return result, totalCount, filteredCount, nil
}

func (r *GormRepository) GetAllFilteredByField(
	model any,
	filters map[string]any,
	preload []string,
	joins []struct {
		Table     string
		Condition string
		JoinType  string
	},
) ([]any, error) {
	query := r.DB.Model(model)

	// Apply preloads
	for _, assoc := range preload {
		query = query.Preload(assoc)
	}

	// Apply joins
	for _, join := range joins {
		if join.JoinType == "" {
			join.JoinType = "JOIN" // Default to INNER JOIN
		}
		query = query.Joins(fmt.Sprintf("%s %s ON %s", join.JoinType, join.Table, join.Condition))
	}

	// Apply filters
	for column, value := range filters {
		if value != nil {
			switch v := value.(type) {
			case []any:
				// Handle slice conditions (e.g., WHERE IN)
				query = query.Where(fmt.Sprintf("%s IN (?)", column), v)
			default:
				// Handle exact match
				query = query.Where(fmt.Sprintf("%s = ?", column), value)
			}
		}
	}

	// Create and populate result slice using reflection
	resultPtr := reflect.New(reflect.SliceOf(reflect.TypeOf(model)))
	if err := query.Find(resultPtr.Interface()).Error; err != nil {
		return nil, fmt.Errorf("error fetching models: %w", err)
	}

	// Convert to []interface{}
	resultSlice := resultPtr.Elem()
	result := make([]any, resultSlice.Len())
	for i := range resultSlice.Len() {
		result[i] = resultSlice.Index(i).Interface()
	}

	return result, nil
}
