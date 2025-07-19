package helper

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/internal/utils/response"
	"gorm.io/gorm"
)

// Create is a generic function to create a new record in the database
// It expects the request body to match the structure of type T
func Create[M any, S any](c *fiber.Ctx, db *gorm.DB) error {

	if c.Request().Body() == nil {
		return response.Fail(c, "EMPTY_BODY", "Request body is required", fiber.StatusBadRequest)
	}

	req := new(S)

	if err := c.BodyParser(req); err != nil {
		return response.Fail(c, "BODY_PARSE_ERROR", "Failed to parse request body", fiber.StatusBadRequest)
	}

	if err := VLD.Struct(req); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}

	var model M

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to process request data", fiber.StatusInternalServerError)
	}

	if err := json.Unmarshal(reqBytes, &model); err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to convert request to model", fiber.StatusInternalServerError)
	}

	if err := db.Create(&model).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to create record: "+err.Error(), fiber.StatusInternalServerError)
	}

	// Return the created record
	return response.OK(c, model, fiber.StatusCreated)
}

// CreateWithValidation is a generic function to create a new record with validation
// It expects a validation function to be passed as a parameter
func CreateWithValidation[M any, S any](c *fiber.Ctx, db *gorm.DB, validateFunc func(M) error) error {

	if c.Request().Body() == nil {
		return response.Fail(c, "EMPTY_BODY", "Request body is required", fiber.StatusBadRequest)
	}

	req := new(S)

	if err := c.BodyParser(req); err != nil {
		return response.Fail(c, "BODY_PARSE_ERROR", "Failed to parse request body", fiber.StatusBadRequest)
	}

	if err := VLD.Struct(req); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}

	var model M

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to process request data", fiber.StatusInternalServerError)
	}

	if err := json.Unmarshal(reqBytes, &model); err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to convert request to model", fiber.StatusInternalServerError)
	}

	// Run custom validation on the model
	if err := validateFunc(model); err != nil {
		return response.Fail(c, "VALIDATION_ERROR", "Validation failed: "+err.Error(), fiber.StatusBadRequest)
	}

	// Create the record in the database
	if err := db.Create(&model).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to create record: "+err.Error(), fiber.StatusInternalServerError)
	}

	// Return the created record
	return response.OK(c, model, fiber.StatusCreated)
}

// CreateWithTransaction is a generic function to create a new record within a database transaction
func CreateWithTransaction[M any, S any](c *fiber.Ctx, db *gorm.DB) error {

	if c.Request().Body() == nil {
		return response.Fail(c, "EMPTY_BODY", "Request body is required", fiber.StatusBadRequest)
	}

	req := new(S)

	if err := c.BodyParser(req); err != nil {
		return response.Fail(c, "BODY_PARSE_ERROR", "Failed to parse request body", fiber.StatusBadRequest)
	}

	if err := VLD.Struct(req); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}

	var model M

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to process request data", fiber.StatusInternalServerError)
	}

	if err := json.Unmarshal(reqBytes, &model); err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to convert request to model", fiber.StatusInternalServerError)
	}

	// Start transaction
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to start transaction: "+err.Error(), fiber.StatusInternalServerError)
	}

	// Create the record in the transaction
	if err := tx.Create(&model).Error; err != nil {
		tx.Rollback()
		return response.Fail(c, "DATABASE_ERROR", "Failed to create record: "+err.Error(), fiber.StatusInternalServerError)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to commit transaction: "+err.Error(), fiber.StatusInternalServerError)
	}

	// Return the created record
	return response.OK(c, model, fiber.StatusCreated)
}
