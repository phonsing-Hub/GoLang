package helper

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/phonsing-Hub/GoLang/internal/utils/response"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UpdateByID[M any, S any](c *fiber.Ctx, db *gorm.DB) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.Fail(c, "INVALID_ID", "Invalid ID format", fiber.StatusBadRequest)
	}

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

	if err := db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Fail(c, "NOT_FOUND", fmt.Sprintf("Record with ID %s not found", idParam), fiber.StatusNotFound)
		}
		return response.Fail(c, "DATABASE_ERROR", "Failed to retrieve record: "+err.Error(), fiber.StatusInternalServerError)
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to process request data", fiber.StatusInternalServerError)
	}

	var updates map[string]interface{}
	if err := json.Unmarshal(reqBytes, &updates); err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to process request data", fiber.StatusInternalServerError)
	}

	protectedFields := []string{"id", "ID", "created_at", "CreatedAt", "updated_at", "UpdatedAt", "deleted_at", "DeletedAt"}
	for _, field := range protectedFields {
		delete(updates, field)
	}

	if len(updates) == 0 {
		return response.Fail(c, "VALIDATION_ERROR", "No valid fields to update", fiber.StatusBadRequest)
	}

	if err := db.Model(&model).Updates(updates).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to update record: "+err.Error(), fiber.StatusInternalServerError)
	}

	if err := db.First(&model, id).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to reload updated record: "+err.Error(), fiber.StatusInternalServerError)
	}

	return response.OK(c, model, fiber.StatusOK)
}

func UpdateByIDWithValidation[M any, S any](c *fiber.Ctx, db *gorm.DB, validateFunc func(M) error) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.Fail(c, "INVALID_ID", "Invalid ID format", fiber.StatusBadRequest)
	}

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

	if err := db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Fail(c, "NOT_FOUND", fmt.Sprintf("Record with ID %s not found", idParam), fiber.StatusNotFound)
		}
		return response.Fail(c, "DATABASE_ERROR", "Failed to retrieve record: "+err.Error(), fiber.StatusInternalServerError)
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to process request data", fiber.StatusInternalServerError)
	}

	var updates map[string]interface{}
	if err := json.Unmarshal(reqBytes, &updates); err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to process request data", fiber.StatusInternalServerError)
	}

	protectedFields := []string{"id", "ID", "created_at", "CreatedAt", "updated_at", "UpdatedAt", "deleted_at", "DeletedAt"}
	for _, field := range protectedFields {
		delete(updates, field)
	}

	if len(updates) == 0 {
		return response.Fail(c, "VALIDATION_ERROR", "No valid fields to update", fiber.StatusBadRequest)
	}

	tx := db.Begin()
	if tx.Error != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to start transaction: "+tx.Error.Error(), fiber.StatusInternalServerError)
	}

	if err := tx.Model(&model).Updates(updates).Error; err != nil {
		tx.Rollback()
		return response.Fail(c, "DATABASE_ERROR", "Failed to apply updates: "+err.Error(), fiber.StatusInternalServerError)
	}

	if err := tx.First(&model, id).Error; err != nil {
		tx.Rollback()
		return response.Fail(c, "DATABASE_ERROR", "Failed to reload updated record: "+err.Error(), fiber.StatusInternalServerError)
	}

	if err := validateFunc(model); err != nil {
		tx.Rollback()
		return response.Fail(c, "VALIDATION_ERROR", "Custom validation failed: "+err.Error(), fiber.StatusBadRequest)
	}

	if err := tx.Commit().Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to commit transaction: "+err.Error(), fiber.StatusInternalServerError)
	}

	return response.OK(c, model, fiber.StatusOK)
}
