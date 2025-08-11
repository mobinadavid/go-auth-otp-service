package scopes

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

const (
	MaxPageSize     = 100
	DefaultPage     = 1
	DefaultPageSize = 10
)

// BuilderModel encapsulates filtering, sorting, and pagination options for database queries
type BuilderModel struct {
	Page          uint
	PageSize      uint
	SortBy        string
	SortOrder     string
	Filters       map[string]interface{}
	Likes         map[string]interface{}
	Joins         []string
	Relations     []string
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	TotalItems    int64
}

// PaginateModel encapsulates the result of a paginated query
type PaginateModel struct {
	PageSize    uint        `json:"page_size"`
	CurrentPage uint        `json:"current_page"`
	TotalPages  int64       `json:"total_pages"`
	TotalItems  int64       `json:"total_items"`
	Items       interface{} `json:"items"`
}

// validatePaginationParams checks and sets default values for pagination parameters
func (bm *BuilderModel) validatePaginationParams() {
	if bm.Page == 0 {
		bm.Page = DefaultPage
	}
	if bm.PageSize > MaxPageSize {
		bm.PageSize = MaxPageSize
	} else if bm.PageSize == 0 {
		bm.PageSize = DefaultPageSize
	}
}

// applySorting applies sorting to the query based on SortBy and SortOrder
func (bm *BuilderModel) applySorting(db *gorm.DB) *gorm.DB {
	if bm.SortBy != "" {
		order := bm.SortBy
		if bm.SortOrder == "desc" {
			order += " desc"
		}
		return db.Order(order)
	}
	return db
}

// applyFilters applies dynamic filtering based on the Filters map in BuilderModel
func (bm *BuilderModel) applyFilters(db *gorm.DB) *gorm.DB {
	for key, value := range bm.Filters {
		parts := strings.Split(key, "_")
		if len(parts) == 2 && parts[1] == "any" {
			column := strings.ReplaceAll(parts[0], ".", "\".\"")
			db = db.Where(clause.Expr{
				SQL:  "? = ANY(\"" + column + "\")",
				Vars: []interface{}{value},
			})
		} else if len(parts) == 3 && parts[2] == "in" {
			column := parts[0] + "_" + parts[1]
			db = db.Where(column+" IN ?", value)
		} else {
			db = db.Where(clause.Expr{SQL: key + " = ?", Vars: []interface{}{value}})
		}
	}
	if bm.CreatedAfter != nil {
		db = db.Where("created_at >= ?", bm.CreatedAfter)
	}
	if bm.CreatedBefore != nil {
		db = db.Where("created_at <= ?", bm.CreatedBefore)
	}
	return db
}

// applyRelations applies relations dynamically
func (bm *BuilderModel) applyRelations(db *gorm.DB) *gorm.DB {
	for _, relation := range bm.Relations {
		db = db.Preload(relation)
	}
	return db
}

// applyFilters applies dynamic filtering based on the Filters map in BuilderModel
func (bm *BuilderModel) applyLikes(db *gorm.DB) *gorm.DB {
	for key, value := range bm.Likes {
		if valStr, ok := value.(string); ok {
			db = db.Where(clause.Expr{SQL: key + " LIKE ?", Vars: []interface{}{"%" + valStr + "%"}})
		}
	}
	return db
}

// applyJoins applies joins to the query
func (bm *BuilderModel) applyJoins(db *gorm.DB) *gorm.DB {
	for _, join := range bm.Joins {
		db = db.Joins(join)
	}
	return db
}

// getTotalItems retrieves the total number of items matching the query
func (bm *BuilderModel) getTotalItems(db *gorm.DB) (int64, error) {
	var totalItems int64
	if err := db.Count(&totalItems).Error; err != nil {
		return 0, err
	}
	return totalItems, nil
}

// calculateOffset calculates the offset for pagination
func (bm *BuilderModel) calculateOffset() int {
	return int((bm.Page - 1) * bm.PageSize)
}

// QueryBuilderScope returns a GORM scope function for paginating results
func (bm *BuilderModel) QueryBuilderScope(db *gorm.DB) (*gorm.DB, error) {
	// Validate and set default pagination parameters
	bm.validatePaginationParams()

	// Apply filters and sorting
	db = bm.applyJoins(db)
	db = bm.applyFilters(db)
	db = bm.applyRelations(db)
	db = bm.applyLikes(db)
	db = bm.applySorting(db)

	// Get total items count
	totalItems, err := bm.getTotalItems(db)
	if err != nil {
		return nil, err
	}
	bm.TotalItems = totalItems

	// Calculate total pages
	totalPages := (totalItems + int64(bm.PageSize) - 1) / int64(bm.PageSize)
	if bm.Page > uint(totalPages) && totalPages > 0 {
		return nil, errors.New("requested page exceeds total number of pages")
	}

	// Calculate offset and apply pagination
	offset := bm.calculateOffset()
	db = db.Offset(offset).Limit(int(bm.PageSize))

	return db, nil
}

// CreatePaginateModel constructs a PaginateModel from the query result
func (bm *BuilderModel) CreatePaginateModel(db *gorm.DB, result interface{}) (*PaginateModel, error) {
	if bm.TotalItems == 0 {
		totalItems, err := bm.getTotalItems(db)
		if err != nil {
			return nil, err
		}
		bm.TotalItems = totalItems
	}

	totalPages := (bm.TotalItems + int64(bm.PageSize) - 1) / int64(bm.PageSize)

	// Execute the query
	if err := db.Find(result).Error; err != nil {
		return nil, err
	}

	return &PaginateModel{
		PageSize:    bm.PageSize,
		CurrentPage: bm.Page,
		TotalPages:  totalPages,
		TotalItems:  bm.TotalItems,
		Items:       result,
	}, nil
}
