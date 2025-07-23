package helper

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/internal/utils/response"
	"gorm.io/gorm"
)

func Create[M any, S any](c *fiber.Ctx, db *gorm.DB, preload ...string) error {
	var schema S
	if err := c.BodyParser(&schema); err != nil {
		return response.Fail(c, "BAD_REQUEST", "Invalid request payload", fiber.StatusBadRequest)
	}

	if err := VLD.Struct(schema); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}

	model := new(M)
	schemaJSON, _ := json.Marshal(schema)
	json.Unmarshal(schemaJSON, model)

	if err := db.Create(model).Error; err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to create record", fiber.StatusInternalServerError)
	}

	tx := db
	for _, p := range preload {
		tx = tx.Preload(p)
	}
	if err := tx.First(model).Error; err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to retrieve created record", fiber.StatusInternalServerError)
	}

	return response.OK(c, model, fiber.StatusCreated)
}

// CreateWithValidation - generic insert with custom validation
func CreateWithValidation[M any, S any](
	c *fiber.Ctx,
	db *gorm.DB,
	validateFn func(c *fiber.Ctx, db *gorm.DB, schema S) error,
	preload ...string,
) error {
	var schema S
	if err := c.BodyParser(&schema); err != nil {
		return response.Fail(c, "BAD_REQUEST", "Invalid request payload", fiber.StatusBadRequest)
	}

	if err := VLD.Struct(schema); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}

	// Custom validation
	if err := validateFn(c, db, schema); err != nil {
		return err
	}

	model := new(M)
	schemaJSON, _ := json.Marshal(schema)
	json.Unmarshal(schemaJSON, model)

	if err := db.Create(model).Error; err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to create record", fiber.StatusInternalServerError)
	}

	tx := db
	for _, p := range preload {
		tx = tx.Preload(p)
	}
	if err := tx.First(model).Error; err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to retrieve created record", fiber.StatusInternalServerError)
	}

	return response.OK(c, model, fiber.StatusCreated)
}
