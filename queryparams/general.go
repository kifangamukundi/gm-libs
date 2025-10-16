package queryparams

import (
	"fmt"

	"slices"

	"github.com/gin-gonic/gin"
)

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func ExtractPaginationParams(c *gin.Context) (int, int, int, string, string, string, any, error) {
	pagination, exists := c.Get("pagination")
	if !exists {
		return 0, 0, 0, "", "", "", nil, fmt.Errorf("pagination info missing")
	}

	paginationMap, ok := pagination.(gin.H)
	if !ok {
		return 0, 0, 0, "", "", "", nil, fmt.Errorf("invalid pagination format")
	}

	page, ok := paginationMap["page"].(int)
	if !ok {
		return 0, 0, 0, "", "", "", nil, fmt.Errorf("invalid or missing page")
	}

	limit, ok := paginationMap["limit"].(int)
	if !ok {
		return 0, 0, 0, "", "", "", nil, fmt.Errorf("invalid or missing limit")
	}

	skip, ok := paginationMap["skip"].(int)
	if !ok {
		return 0, 0, 0, "", "", "", nil, fmt.Errorf("invalid or missing skip")
	}

	sortOrder, exists := c.Get("sortOrder")
	if !exists {
		return 0, 0, 0, "", "", "", nil, fmt.Errorf("sort order missing")
	}

	sortByColumn, exists := c.Get("sortByColumn")
	if !exists {
		return 0, 0, 0, "", "", "", nil, fmt.Errorf("sort by column missing")
	}

	searchRegex, exists := c.Get("searchRegex")
	if !exists {
		searchRegex = ""
	}

	filterCriteria, exists := c.Get("filterCriteria")
	if !exists {
		filterCriteria = nil
	}

	return page, limit, skip, sortOrder.(string), sortByColumn.(string), searchRegex.(string), filterCriteria, nil
}
