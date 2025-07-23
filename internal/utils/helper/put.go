package helper

import (
	"strconv"
	"time"

	"github.com/phonsing-Hub/GoLang/internal/utils/response"
	"github.com/phonsing-Hub/GoLang/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UpdateByID[M any, S any](c *fiber.Ctx, db *gorm.DB, preload ...string) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.Fail(c, "INVALID_ID", "Invalid ID format", fiber.StatusBadRequest)
	}

	var schema S
	if err := c.BodyParser(&schema); err != nil {
		return response.Fail(c, "BAD_REQUEST", "Invalid request payload", fiber.StatusBadRequest)
	}

	if err := VLD.Struct(schema); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}

	model := new(M)
	if err := db.First(model, id).Error; err != nil {
		return response.Fail(c, "NOT_FOUND", "data not found", fiber.StatusNotFound)
	}

	if err := db.Model(model).Updates(schema).Error; err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to update record", fiber.StatusInternalServerError)
	}

	if err := db.Model(model).Update("updated_at", time.Now()).Error; err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to update record", fiber.StatusInternalServerError)
	}

	tx := db
	for _, p := range preload {
		tx = tx.Preload(p)
	}
	if err := tx.First(model, id).Error; err != nil {
		return response.Fail(c, "NOT_FOUND", "data not found", fiber.StatusNotFound)
	}

	return response.OK(c, model, fiber.StatusOK)
}

func UpdateByClaims[M any, S any](c *fiber.Ctx, db *gorm.DB, preload ...string) error {
	user := c.Locals("user").(*jwt.Claims)
	if user == nil {
		return response.Fail(c, "UNAUTHORIZED", "Unauthorized", fiber.StatusUnauthorized)
	}
	id := user.UserID

	var schema S
	if err := c.BodyParser(&schema); err != nil {
		return response.Fail(c, "BAD_REQUEST", "Invalid request payload", fiber.StatusBadRequest)
	}

	if err := VLD.Struct(schema); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}

	model := new(M)
	if err := db.First(model, id).Error; err != nil {
		return response.Fail(c, "NOT_FOUND", "data not found", fiber.StatusNotFound)
	}

	if err := db.Model(model).Updates(schema).Error; err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to update record", fiber.StatusInternalServerError)
	}

	if err := db.Model(model).Update("updated_at", time.Now()).Error; err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to update record", fiber.StatusInternalServerError)
	}

	tx := db
	for _, p := range preload {
		tx = tx.Preload(p)
	}
	if err := tx.First(model, id).Error; err != nil {
		return response.Fail(c, "NOT_FOUND", "data not found", fiber.StatusNotFound)
	}

	return response.OK(c, model, fiber.StatusOK)
}

func UpdateWithValidation[M any, S any](
	c *fiber.Ctx,
	db *gorm.DB,
	validateFn func(c *fiber.Ctx, db *gorm.DB, args ...any) error,
	where string,
	whereArgs []any,
	preload ...string,
) error {
	if err := validateFn(c, db, whereArgs...); err != nil {
		return err
	}

	var schema S
	if err := c.BodyParser(&schema); err != nil {
		return response.Fail(c, "BAD_REQUEST", "Invalid request payload", fiber.StatusBadRequest)
	}

	if err := VLD.Struct(schema); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}

	if err := db.Model(new(M)).Where(where, whereArgs...).Updates(schema).Error; err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to update record", fiber.StatusInternalServerError)
	}

		if err := db.Model(new(M)).Where(where, whereArgs...).Update("updated_at", time.Now()).Error; err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to update record", fiber.StatusInternalServerError)
	}

	model := new(M)
	tx := db
	for _, p := range preload {
		tx = tx.Preload(p)
	}
	if err := tx.First(model, whereArgs...).Error; err != nil {
		return response.Fail(c, "NOT_FOUND", "data not found", fiber.StatusNotFound)
	}

	return response.OK(c, model, fiber.StatusOK)
}
