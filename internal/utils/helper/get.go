package helper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/internal/utils/response"
	"gorm.io/gorm"
)

type PaginatedResponse[T any] struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Data  []T   `json:"data"`
}

// FindAllWithPreload is similar to FindAll but supports preloading relationships
func FindAllWithPreload[T any](c *fiber.Ctx, db *gorm.DB, preloads ...string) error {
	var models []T
	var total int64

	query := db.Model(new(T))

	// Apply preloads
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	queryParams := c.Queries()

	for key, value := range queryParams {
		values := strings.Split(value, ",")

		switch {
		case strings.HasPrefix(key, "search[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "search["), "]")
			query = query.Where(fmt.Sprintf("%s ILIKE ?", colName), "%"+value+"%")

		case strings.HasPrefix(key, "search_cols[") && strings.HasSuffix(key, "]"):
			colNames := strings.Split(strings.TrimSuffix(strings.TrimPrefix(key, "search_cols["), "]"), "|")
			var orConditions []string
			var orArgs []interface{}
			for _, col := range colNames {
				orConditions = append(orConditions, fmt.Sprintf("%s ILIKE ?", col))
				orArgs = append(orArgs, "%"+value+"%")
			}
			query = query.Where(strings.Join(orConditions, " OR "), orArgs...)

		case strings.HasPrefix(key, "filter_not[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "filter_not["), "]")
			query = query.Where(fmt.Sprintf("%s NOT IN (?)", colName), values)

		case strings.HasPrefix(key, "filterrange[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "filterrange["), "]")
			rangeValues := strings.Split(value, "|")
			if len(rangeValues) == 2 {
				if rangeValues[0] != "-" {
					query = query.Where(fmt.Sprintf("%s >= ?", colName), rangeValues[0])
				}
				if rangeValues[1] != "-" {
					query = query.Where(fmt.Sprintf("%s <= ?", colName), rangeValues[1])
				}
			}

		case key != "page" && key != "limit" && key != "sort_by" && key != "sort_order":
			if len(values) == 1 {
				if strings.ToLower(values[0]) == "null" {
					query = query.Where(fmt.Sprintf("%s IS NULL", key))
				} else {
					query = query.Where(fmt.Sprintf("%s = ?", key), values[0])
				}
			} else {
				query = query.Where(fmt.Sprintf("%s IN (?)", key), values)
			}
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to count records: "+err.Error(), fiber.StatusInternalServerError)
	}

	sortBy := c.Query("sort_by", "id")
	sortOrder := c.Query("sort_order", "asc")
	if strings.ToLower(sortOrder) == "desc" {
		query = query.Order(fmt.Sprintf("%s desc", sortBy))
	} else {
		query = query.Order(fmt.Sprintf("%s asc", sortBy))
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&models).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to fetch records: "+err.Error(), fiber.StatusInternalServerError)
	}

	return response.OK(c, &PaginatedResponse[T]{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  models,
	})
}

// FindAllWithJoins supports custom joins and preloads
func FindAllWithJoins[T any](c *fiber.Ctx, db *gorm.DB, joins []string, preloads []string) error {
	var models []T
	var total int64

	query := db.Model(new(T))

	// Apply joins
	for _, join := range joins {
		query = query.Joins(join)
	}

	// Apply preloads
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	queryParams := c.Queries()

	for key, value := range queryParams {
		values := strings.Split(value, ",")

		switch {
		case strings.HasPrefix(key, "search[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "search["), "]")
			query = query.Where(fmt.Sprintf("%s ILIKE ?", colName), "%"+value+"%")

		case strings.HasPrefix(key, "search_cols[") && strings.HasSuffix(key, "]"):
			colNames := strings.Split(strings.TrimSuffix(strings.TrimPrefix(key, "search_cols["), "]"), "|")
			var orConditions []string
			var orArgs []interface{}
			for _, col := range colNames {
				orConditions = append(orConditions, fmt.Sprintf("%s ILIKE ?", col))
				orArgs = append(orArgs, "%"+value+"%")
			}
			query = query.Where(strings.Join(orConditions, " OR "), orArgs...)

		case strings.HasPrefix(key, "filter_not[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "filter_not["), "]")
			query = query.Where(fmt.Sprintf("%s NOT IN (?)", colName), values)

		case strings.HasPrefix(key, "filterrange[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "filterrange["), "]")
			rangeValues := strings.Split(value, "|")
			if len(rangeValues) == 2 {
				if rangeValues[0] != "-" {
					query = query.Where(fmt.Sprintf("%s >= ?", colName), rangeValues[0])
				}
				if rangeValues[1] != "-" {
					query = query.Where(fmt.Sprintf("%s <= ?", colName), rangeValues[1])
				}
			}

		case key != "page" && key != "limit" && key != "sort_by" && key != "sort_order":
			if len(values) == 1 {
				if strings.ToLower(values[0]) == "null" {
					query = query.Where(fmt.Sprintf("%s IS NULL", key))
				} else {
					query = query.Where(fmt.Sprintf("%s = ?", key), values[0])
				}
			} else {
				query = query.Where(fmt.Sprintf("%s IN (?)", key), values)
			}
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to count records: "+err.Error(), fiber.StatusInternalServerError)
	}

	sortBy := c.Query("sort_by", "id")
	sortOrder := c.Query("sort_order", "asc")
	if strings.ToLower(sortOrder) == "desc" {
		query = query.Order(fmt.Sprintf("%s desc", sortBy))
	} else {
		query = query.Order(fmt.Sprintf("%s asc", sortBy))
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&models).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to fetch records: "+err.Error(), fiber.StatusInternalServerError)
	}

	return response.OK(c, &PaginatedResponse[T]{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  models,
	})
}

func FindAll[T any](c *fiber.Ctx, db *gorm.DB) error {
	var models []T
	var total int64

	query := db.Model(new(T))

	queryParams := c.Queries()

	for key, value := range queryParams {
		values := strings.Split(value, ",")

		switch {
		case strings.HasPrefix(key, "search[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "search["), "]")
			query = query.Where(fmt.Sprintf("%s ILIKE ?", colName), "%"+value+"%")

		case strings.HasPrefix(key, "search_cols[") && strings.HasSuffix(key, "]"):
			colNames := strings.Split(strings.TrimSuffix(strings.TrimPrefix(key, "search_cols["), "]"), "|")
			var orConditions []string
			var orArgs []interface{}
			for _, col := range colNames {
				orConditions = append(orConditions, fmt.Sprintf("%s ILIKE ?", col))
				orArgs = append(orArgs, "%"+value+"%")
			}
			query = query.Where(strings.Join(orConditions, " OR "), orArgs...)

		case strings.HasPrefix(key, "filter_not[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "filter_not["), "]")
			query = query.Where(fmt.Sprintf("%s NOT IN (?)", colName), values)

		case strings.HasPrefix(key, "filterrange[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "filterrange["), "]")
			rangeValues := strings.Split(value, "|")
			if len(rangeValues) == 2 {
				if rangeValues[0] != "-" {
					query = query.Where(fmt.Sprintf("%s >= ?", colName), rangeValues[0])
				}
				if rangeValues[1] != "-" {
					query = query.Where(fmt.Sprintf("%s <= ?", colName), rangeValues[1])
				}
			}

		case key != "page" && key != "limit" && key != "sort_by" && key != "sort_order":
			if len(values) == 1 {
				if strings.ToLower(values[0]) == "null" {
					query = query.Where(fmt.Sprintf("%s IS NULL", key))
				} else {
					query = query.Where(fmt.Sprintf("%s = ?", key), values[0])
				}
			} else {
				query = query.Where(fmt.Sprintf("%s IN (?)", key), values)
			}
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to count records: "+err.Error(), fiber.StatusInternalServerError)
	}

	sortBy := c.Query("sort_by", "id")
	sortOrder := c.Query("sort_order", "asc")
	if strings.ToLower(sortOrder) == "desc" {
		query = query.Order(fmt.Sprintf("%s desc", sortBy))
	} else {
		query = query.Order(fmt.Sprintf("%s asc", sortBy))
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&models).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to fetch records: "+err.Error(), fiber.StatusInternalServerError)
	}

	return response.OK(c, &PaginatedResponse[T]{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  models,
	})
}

// FindByIDWithPreload supports preloading relationships for single record
func FindByIDWithPreload[T any](c *fiber.Ctx, db *gorm.DB, preloads ...string) error {
	idParam := c.Params("id")

	// Validate ID format
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.Fail(c, "INVALID_ID", "Invalid ID format", fiber.StatusBadRequest)
	}

	query := db.Model(new(T))

	// Apply preloads
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	var model T
	if err := query.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Fail(c, "NOT_FOUND", fmt.Sprintf("Record with ID %s not found", idParam), fiber.StatusNotFound)
		}
		return response.Fail(c, "DATABASE_ERROR", "Failed to retrieve record: "+err.Error(), fiber.StatusInternalServerError)
	}
	return response.OK(c, model)
}

func FindByID[T any](c *fiber.Ctx, db *gorm.DB) error {
	idParam := c.Params("id")

	// Validate ID format
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.Fail(c, "INVALID_ID", "Invalid ID format", fiber.StatusBadRequest)
	}

	var model T
	if err := db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Fail(c, "NOT_FOUND", fmt.Sprintf("Record with ID %s not found", idParam), fiber.StatusNotFound)
		}
		return response.Fail(c, "DATABASE_ERROR", "Failed to retrieve record: "+err.Error(), fiber.StatusInternalServerError)
	}
	return response.OK(c, model)
}
