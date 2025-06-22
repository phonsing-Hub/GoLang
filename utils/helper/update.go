package helper

import (
	"fmt"
	"github.com/phonsing-Hub/GoLang/utils/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UpdateByID[T any](c *fiber.Ctx, db *gorm.DB) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.Fail(c, "INVALID_ID", "Invalid ID format", fiber.StatusBadRequest)
	}

	var model T

	if err := db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Fail(c, "NOT_FOUND", fmt.Sprintf("Record with ID %s not found", idParam), fiber.StatusNotFound)
		}
		return response.Fail(c, "DATABASE_ERROR", fmt.Sprintf("Failed to retrieve record: %v", err), fiber.StatusInternalServerError)
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return response.Fail(c, "BODY_PARSE_ERROR", "Failed to parse request body", fiber.StatusBadRequest)
	}

	delete(updates, "id")
	delete(updates, "ID")
	delete(updates, "created_at")
	delete(updates, "CreatedAt")
	delete(updates, "updated_at")
	delete(updates, "UpdatedAt")
	delete(updates, "deleted_at")
	delete(updates, "DeletedAt")

	if len(updates) == 0 {
		return response.Fail(c, "VALIDATION_ERROR", "No valid fields to update", fiber.StatusBadRequest)
	}

	if err := db.Model(&model).Updates(&updates).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", fmt.Sprintf("Failed to update record: %v", err), fiber.StatusInternalServerError)
	}

	if err := db.First(&model, id).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", fmt.Sprintf("Failed to reload updated record: %v", err), fiber.StatusInternalServerError)
	}

	return response.OK(c, model, fiber.StatusOK)
}
