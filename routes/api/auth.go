package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/database"
	"github.com/phonsing-Hub/GoLang/database/models"
	"github.com/phonsing-Hub/GoLang/database/schema"
	"github.com/phonsing-Hub/GoLang/middleware"
	"github.com/phonsing-Hub/GoLang/pkg/auth"
	"github.com/phonsing-Hub/GoLang/pkg/jwt"
	"github.com/phonsing-Hub/GoLang/utils/helper"
	"github.com/phonsing-Hub/GoLang/utils/response"
	"gorm.io/gorm"
)

func SetupAuthRoutes(router *fiber.App) {
	authGroup := router.Group("/auth")
	authGroup.Post("/login", login)
	authGroup.Get("/profile", middleware.JWTAuthMiddleware(), profile)
}

func login(c *fiber.Ctx) error {
	req := new(schema.LoginRequest)

	if err := c.BodyParser(req); err != nil {
		return response.Fail(c, "BODY_NOT_FOUND", "body not found", fiber.StatusOK)
	}

	if err := helper.VLD.Struct(req); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}
	var existingUser models.UserCredentials
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err != nil {
		return response.Fail(c, "USER_NOT_FOUND", "user not found", fiber.StatusNotFound)
	}

	if !auth.CheckPasswordHash(req.Password, existingUser.PasswordHash) {
		return response.Fail(c, "INVALID_CREDENTIALS", "invalid credentials", fiber.StatusUnauthorized)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(existingUser.UserID, existingUser.Username)
	if err != nil {
		return response.Fail(c, "TOKEN_GENERATION_FAILED", "failed to generate token", fiber.StatusInternalServerError)
	}

	return response.OK(c, fiber.Map{
		"token": token,
	}, fiber.StatusOK)
}

func profile(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*jwt.Claims)
	var user models.Users
	err := database.DB.
		Joins("CurrentStatus").
		Preload("Locations").  
		Where("users.id = ?", userClaims.UserID).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Fail(c, "NOT_FOUND", "User not found", fiber.StatusUnauthorized)
		}
		return response.Fail(c, "DATABASE_ERROR", "Failed to retrieve user data: "+err.Error(), fiber.StatusInternalServerError)
	}

	return response.OK(c, user, fiber.StatusOK)
}
