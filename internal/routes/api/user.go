package api

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/internal/database"
	"github.com/phonsing-Hub/GoLang/internal/database/models"

	"github.com/phonsing-Hub/GoLang/internal/database/schema"
	"github.com/phonsing-Hub/GoLang/internal/middleware"
	"github.com/phonsing-Hub/GoLang/internal/utils/helper"
	"github.com/phonsing-Hub/GoLang/internal/utils/response"
	"github.com/phonsing-Hub/GoLang/pkg/jwt"
)

func SetupUserhRoutes(router *fiber.App) {
	userGroup := router.Group("/profile")
	userGroup.Use(middleware.JWTAuthMiddleware())
	userGroup.Get("", profileHandler)
	userGroup.Put("", updateProfileHandler)
	userGroup.Put("/avatar", updateAvatarHandler)
}

// profileHandler retrieves the user's profile information
// @Summary Get User Profile
// @Description Retrieve the authenticated user's profile information
// @Tags PROFILE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SWSuccessResponse
// @Failure 401 {object} response.SWErrorResponse
// @Failure 500 {object} response.SWErrorResponse
// @Router /profile [get]
func profileHandler(c *fiber.Ctx) error {
	db := database.DB
	user := c.Locals("user").(*jwt.Claims)
	if user == nil {
		return response.Fail(c, "UNAUTHORIZED", "Unauthorized", fiber.StatusUnauthorized)
	}
	dbUser := &models.User{}
	if err := db.Preload("Status").First(dbUser, user.UserID).Error; err != nil {
		return response.Fail(c, "NOT_FOUND", "User not found", fiber.StatusNotFound)
	}
	return response.OK(c, dbUser, fiber.StatusOK)
}

// updateProfileHandler updates the user's profile information
// @Summary Update User Profile
// @Description Update the authenticated user's profile information
// @Tags PROFILE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body schema.UpdateUser true "User Profile Data"
// @Success 200 {object} response.SWSuccessResponse
// @Failure 400 {object} response.SWErrorResponse
// @Failure 401 {object} response.SWErrorResponse
// @Failure 500 {object} response.SWErrorResponse
// @Router /profile [put]
func updateProfileHandler(c *fiber.Ctx) error {
	return helper.UpdateByClaims[models.User, schema.UpdateUser](c, database.DB, "Status")
}

// updateAvatarHandler updates the user's avatar
// @Summary Update User Avatar
// @Description Update the authenticated user's avatar
// @Tags PROFILE
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param avatar formData file true "Avatar image file (PNG, JPEG, WEBP, max 2MB)"
// @Success 200 {object} response.SWSuccessResponse
// @Failure 400 {object} response.SWErrorResponse
// @Failure 401 {object} response.SWErrorResponse
// @Failure 500 {object} response.SWErrorResponse
// @Router /profile/avatar [put]
func updateAvatarHandler(c *fiber.Ctx) error {
	db := database.DB
	user := c.Locals("user").(*jwt.Claims)
	if user == nil {
		return response.Fail(c, "UNAUTHORIZED", "Unauthorized", fiber.StatusUnauthorized)
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return response.Fail(c, "BAD_REQUEST", "Avatar file is required", fiber.StatusBadRequest)
	}

	currentUser := &models.User{}
	if err := db.First(currentUser, user.UserID).Error; err != nil {
		return response.Fail(c, "NOT_FOUND", "User not found", fiber.StatusNotFound)
	}

	if currentUser.Avatar != "" {
		oldAvatarPath := filepath.Join("./static/uploads", currentUser.Avatar)
		if _, err := os.Stat(oldAvatarPath); err == nil {
			os.Remove(oldAvatarPath) 
		}
	}

	newFileName, err := helper.ValidateAndRenameAvatar(file)
	if err != nil {
		return response.Fail(c, "BAD_REQUEST", err.Error(), fiber.StatusBadRequest)
	}

	if err := helper.PrepareAvatarPath(newFileName); err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to prepare upload directory", fiber.StatusInternalServerError)
	}

	savePath := "./static/uploads/" + newFileName
	if err := c.SaveFile(file, savePath); err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to save avatar", fiber.StatusInternalServerError)
	}

	currentUser.Avatar = newFileName
	if err := db.Save(&currentUser).Error; err != nil {
		return response.Fail(c, "INTERNAL_SERVER_ERROR", "Failed to update avatar", fiber.StatusInternalServerError)
	}
	
	return response.OK(c, fiber.Map{
		"avatar": newFileName,
	}, fiber.StatusOK)
}
