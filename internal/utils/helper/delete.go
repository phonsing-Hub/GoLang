package helper

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/internal/utils/response"
	"gorm.io/gorm"
)

// DeleteByID is a generic function to delete a record by ID
func DeleteByID[T any](c *fiber.Ctx, db *gorm.DB) error {
	idParam := c.Params("id")

	// Validate ID format
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.Fail(c, "INVALID_ID", "Invalid ID format", fiber.StatusBadRequest)
	}

	var model T

	// Check if record exists
	if err := db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Fail(c, "NOT_FOUND", fmt.Sprintf("Record with ID %s not found", idParam), fiber.StatusNotFound)
		}
		return response.Fail(c, "DATABASE_ERROR", "Failed to retrieve record: "+err.Error(), fiber.StatusInternalServerError)
	}

	// Delete the record
	if err := db.Unscoped().Delete(&model).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to delete record: "+err.Error(), fiber.StatusInternalServerError)
	}

	return response.OK(c, fiber.Map{
		"message": "Record deleted successfully",
		"id":      id,
	}, fiber.StatusOK)
}

// DeleteByIDWithValidation is a generic function to delete a record with custom validation
func DeleteByIDWithValidation[T any](c *fiber.Ctx, db *gorm.DB, validateFunc func(T) error) error {
	idParam := c.Params("id")

	// Validate ID format
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.Fail(c, "INVALID_ID", "Invalid ID format", fiber.StatusBadRequest)
	}

	var model T

	// Check if record exists
	if err := db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Fail(c, "NOT_FOUND", fmt.Sprintf("Record with ID %s not found", idParam), fiber.StatusNotFound)
		}
		return response.Fail(c, "DATABASE_ERROR", "Failed to retrieve record: "+err.Error(), fiber.StatusInternalServerError)
	}

	// Run custom validation
	if err := validateFunc(model); err != nil {
		return response.Fail(c, "VALIDATION_ERROR", "Validation failed: "+err.Error(), fiber.StatusBadRequest)
	}

	// Delete the record
	if err := db.Delete(&model).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to delete record: "+err.Error(), fiber.StatusInternalServerError)
	}

	return response.OK(c, fiber.Map{
		"message": "Record deleted successfully",
		"id":      id,
	}, fiber.StatusOK)
}

// DeleteByIDWithTransaction is a generic function to delete a record within a database transaction
func DeleteByIDWithTransaction[T any](c *fiber.Ctx, db *gorm.DB) error {
	idParam := c.Params("id")

	// Validate ID format
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.Fail(c, "INVALID_ID", "Invalid ID format", fiber.StatusBadRequest)
	}

	var model T

	// Start transaction
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to start transaction: "+err.Error(), fiber.StatusInternalServerError)
	}

	// Check if record exists
	if err := tx.First(&model, id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return response.Fail(c, "NOT_FOUND", fmt.Sprintf("Record with ID %s not found", idParam), fiber.StatusNotFound)
		}
		return response.Fail(c, "DATABASE_ERROR", "Failed to retrieve record: "+err.Error(), fiber.StatusInternalServerError)
	}

	// Delete the record
	if err := tx.Delete(&model).Error; err != nil {
		tx.Rollback()
		return response.Fail(c, "DATABASE_ERROR", "Failed to delete record: "+err.Error(), fiber.StatusInternalServerError)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to commit transaction: "+err.Error(), fiber.StatusInternalServerError)
	}

	return response.OK(c, fiber.Map{
		"message": "Record deleted successfully",
		"id":      id,
	}, fiber.StatusOK)
}

// SoftDeleteByID is a generic function to soft delete a record (set deleted_at)
func SoftDeleteByID[T any](c *fiber.Ctx, db *gorm.DB) error {
	idParam := c.Params("id")

	// Validate ID format
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.Fail(c, "INVALID_ID", "Invalid ID format", fiber.StatusBadRequest)
	}

	var model T

	// Check if record exists
	if err := db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Fail(c, "NOT_FOUND", fmt.Sprintf("Record with ID %s not found", idParam), fiber.StatusNotFound)
		}
		return response.Fail(c, "DATABASE_ERROR", "Failed to retrieve record: "+err.Error(), fiber.StatusInternalServerError)
	}

	// Soft delete the record
	if err := db.Delete(&model).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to soft delete record: "+err.Error(), fiber.StatusInternalServerError)
	}

	return response.OK(c, fiber.Map{
		"message": "Record soft deleted successfully",
		"id":      id,
	}, fiber.StatusOK)
}
