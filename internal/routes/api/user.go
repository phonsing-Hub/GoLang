package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/internal/database"
	"github.com/phonsing-Hub/GoLang/internal/database/models"
	"github.com/phonsing-Hub/GoLang/internal/database/schema"
	"github.com/phonsing-Hub/GoLang/internal/middleware"
	"github.com/phonsing-Hub/GoLang/internal/utils/helper"
	"github.com/phonsing-Hub/GoLang/internal/utils/response"
	"github.com/phonsing-Hub/GoLang/pkg/auth"
	"os"
	"path/filepath"
)

func SetupUserRoutes(router *fiber.App) {
	userGroup := router.Group("/users")
	userGroup.Use(middleware.JWTAuthMiddleware())

	userGroup.Get("/", get_users)
	userGroup.Get("/:id", get_user_id)
	userGroup.Post("/", create_user)
	userGroup.Put("/:id", update_user)

	userGroup.Post("/:id/avatar", upload_user_avatar)
	userGroup.Delete("/:id/avatar", delete_user_avatar)
}

func get_users(c *fiber.Ctx) error {
	return helper.FindAll[models.Users](c, database.DB)
}

func get_user_id(c *fiber.Ctx) error {
	return helper.FindByID[models.Users](c, database.DB)
}

func create_user(c *fiber.Ctx) error {
	req := new(schema.CreateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return response.Fail(c, "BODY_NOT_FOUND", "body not found", fiber.StatusOK)
	}

	if err := helper.VLD.Struct(req); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}

	nwedateOfBirth, _ := helper.StringToDate(req.DateOfBirth, "2006-01-02")
	hastPassword, _ := auth.HashPassword(req.Password)

	new_user := models.Users{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DisplayName: req.DisplayName,
		Gender:      req.Gender,
		DateOfBirth: &nwedateOfBirth,
	}

	tx := database.DB.Begin()
	if err := tx.Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", err.Error(), fiber.StatusInternalServerError)
	}
	if err := tx.Create(&new_user).Error; err != nil {
		tx.Rollback()
		return response.Fail(c, "DATABASE_ERROR", err.Error(), fiber.StatusInternalServerError)
	}

	newCredential := models.UserCredentials{
		UserID:         new_user.ID,
		CredentialType: "password",
		Username:       req.Email,
		PasswordHash:   hastPassword,
	}
	if err := tx.Create(&newCredential).Error; err != nil {
		tx.Rollback()
		return response.Fail(c, "DATABASE_ERROR", "Failed to create user credentials", fiber.StatusInternalServerError)
	}

	if req.Location != nil {
		newLocation := models.UserLocation{
			UserID:       new_user.ID,
			LocationType: req.Location.LocationType,
			AddressLine1: req.Location.AddressLine1,
			AddressLine2: req.Location.AddressLine2,
			City:         req.Location.City,
			PostalCode:   req.Location.PostalCode,
			Country:      req.Location.Country,
			IsDefault:    true, // Mark as default if it's the first location
		}
		if err := tx.Create(&newLocation).Error; err != nil {
			tx.Rollback()
			return response.Fail(c, "DATABASE_ERROR", "Failed to create user location", fiber.StatusInternalServerError)
		}
	}

	tx.Commit()
	return response.OK(c, fiber.Map{
		"user_id": new_user.ID,
		"email":   new_user.Email,
	}, fiber.StatusCreated)

}

func update_user(c *fiber.Ctx) error {
	return helper.UpdateByID[models.Users](c, database.DB)
}

func upload_user_avatar(c *fiber.Ctx) error {
	id := c.Params("id")

	file, err := c.FormFile("avatar")
	if err != nil {
		return response.Fail(c, "BAD_REQUEST", "Avatar is required", fiber.StatusBadRequest)
	}

	var user models.Users
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		return response.Fail(c, "NOT_FOUND", "User not found", fiber.StatusNotFound)
	}

	if user.Avatar != "" {
		oldPath := filepath.Join("./static/uploads", user.Avatar)
		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			return response.Fail(c, "INTERNAL_ERROR", "Failed to delete old avatar", fiber.StatusInternalServerError)
		}
	}

	savePath, err := helper.ValidateAndRenameAvatar(file)
	if err != nil {
		return response.Fail(c, "BAD_REQUEST", err.Error(), fiber.StatusBadRequest)
	}

	if err := helper.PrepareAvatarPath(savePath); err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to create directory", fiber.StatusInternalServerError)
	}

	fullPath := fmt.Sprintf("./static/uploads/%s", savePath)

	if err := c.SaveFile(file, fullPath); err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to save avatar", fiber.StatusInternalServerError)
	}

	if err := database.DB.Model(&user).Update("avatar", savePath).Error; err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to update avatar in database", fiber.StatusInternalServerError)
	}

	return response.OK(c, fiber.Map{
		"message":   "Avatar uploaded successfully",
		"avatarUrl": fmt.Sprintf("/static/%s", savePath),
	}, fiber.StatusCreated)
}

func delete_user_avatar(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.Users
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		return response.Fail(c, "NOT_FOUND", "User not found", fiber.StatusNotFound)
	}

	if user.Avatar == "" {
		return response.Fail(c, "BAD_REQUEST", "User has no avatar", fiber.StatusBadRequest)
	}

	fullPath := fmt.Sprintf("./static/uploads/%s", user.Avatar)

	if err := os.Remove(fullPath); err != nil {
		if !os.IsNotExist(err) {
			return response.Fail(c, "INTERNAL_ERROR", "Failed to delete avatar file", fiber.StatusInternalServerError)
		}
	}

	tx := database.DB.Begin()
	if err := tx.Model(&user).Update("avatar", "").Error; err != nil {
		tx.Rollback()
		return response.Fail(c, "INTERNAL_ERROR", "Failed to clear avatar in database", fiber.StatusInternalServerError)
	}
	tx.Commit()

	return response.OK(c, fiber.Map{
		"message": "Avatar deleted successfully",
	}, fiber.StatusOK)
}
