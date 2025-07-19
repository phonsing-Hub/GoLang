package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/internal/database"
	"github.com/phonsing-Hub/GoLang/internal/database/models"
	"github.com/phonsing-Hub/GoLang/internal/database/schema"
	"github.com/phonsing-Hub/GoLang/internal/utils/helper"
	"github.com/phonsing-Hub/GoLang/internal/utils/response"
	"github.com/phonsing-Hub/GoLang/internal/middleware"
	"github.com/phonsing-Hub/GoLang/pkg/auth"
	"github.com/phonsing-Hub/GoLang/pkg/jwt"
	"time"
)

func SetupAuthRoutes(router *fiber.App) {
	authGroup := router.Group("/auth")
	authGroup.Post("/register", registerUser)
	authGroup.Post("/login", loginUser)
	authGroup.Get("/userinfo", middleware.JWTAuthMiddleware(), getUserInfo)
}

// registerUser handles user registration
// @Summary User Registration
// @Description Register a new user with email, username, and password
// @Tags AUTH
// @Accept json
// @Produce json
// @Param register body schema.RegisterRequest true "Register Request"
// @Success 201 {object} response.SWSuccessResponse
// @Failure 400 {object} response.SWErrorResponse
// @Failure 500 {object} response.SWErrorResponse
// @Router /auth/register [post]
func registerUser(c *fiber.Ctx) error {
	// Parse request
	var req schema.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, "BODY_PARSE_ERROR", "Failed to parse request body", fiber.StatusBadRequest)
	}

	// Validate request
	if err := helper.VLD.Struct(req); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}

	// Check if email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return response.Fail(c, "EMAIL_EXISTS", "Email already registered", fiber.StatusBadRequest)
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to hash password", fiber.StatusInternalServerError)
	}

	// Begin transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user
	now := time.Now()
	user := models.User{
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DisplayName: req.DisplayName,
		LastLoginAt: &now,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return response.Fail(c, "DATABASE_ERROR", "Failed to create user", fiber.StatusInternalServerError)
	}

	// Create auth method
	authMethod := models.UserAuthMethod{
		UserID:       user.ID,
		AuthType:     "password",
		IsPrimary:    true,
		PasswordHash: hashedPassword,
	}

	if err := tx.Create(&authMethod).Error; err != nil {
		tx.Rollback()
		return response.Fail(c, "DATABASE_ERROR", "Failed to create auth method", fiber.StatusInternalServerError)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to commit transaction", fiber.StatusInternalServerError)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to generate token", fiber.StatusInternalServerError)
	}

	userInfo := schema.UserInfo{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		DisplayName:     user.DisplayName,
		Avatar:          user.Avatar,
		IsEmailVerified: user.IsEmailVerified,
		LastLoginAt:     user.LastLoginAt,
	}

	return response.OK(c, fiber.Map{
		"user":  userInfo,
		"token": token,
	}, fiber.StatusCreated)
}

// loginUser handles user login
// @Summary User Login
// @Description Login user with email and password
// @Tags AUTH
// @Accept json
// @Produce json
// @Param login body schema.LoginRequest true "Login Request"
// @Success 200 {object} response.SWSuccessResponse
// @Failure 400 {object} response.SWErrorResponse
// @Failure 401 {object} response.SWErrorResponse
// @Failure 500 {object} response.SWErrorResponse
// @Router /auth/login [post]
func loginUser(c *fiber.Ctx) error {
	// Parse request
	var req schema.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, "BODY_PARSE_ERROR", "Failed to parse request body", fiber.StatusBadRequest)
	}

	// Validate request
	if err := helper.VLD.Struct(req); err != nil {
		return response.Fail(c, "VALIDATION_FAILED", err.Error(), fiber.StatusBadRequest)
	}

	// Find user by email or username
	var user models.User
	if err := database.DB.Preload("AuthMethods").Where("email = ? ", req.Email).First(&user).Error; err != nil {
		return response.Fail(c, "USER_NOT_FOUND", "Invalid email or password", fiber.StatusNotFound)
	}

	// Verify password
	if !auth.CheckPasswordHash(req.Password, user.AuthMethods[0].PasswordHash) {
		return response.Fail(c, "INVALID_CREDENTIALS", "Invalid email or password", fiber.StatusUnauthorized)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to generate token", fiber.StatusInternalServerError)
	}

	now := time.Now()

	if err := database.DB.Model(&user).Update("last_login_at", &now).Error; err != nil {
    return response.Fail(c, "DB_ERROR", "Failed to update last login", fiber.StatusInternalServerError)
}
	userInfo := schema.UserInfo{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		DisplayName:     user.DisplayName,
		Avatar:          user.Avatar,
		IsEmailVerified: user.IsEmailVerified,
		LastLoginAt:     user.LastLoginAt,
	}

	return response.OK(c, fiber.Map{
		"user":  userInfo,
		"token": token,
	}, fiber.StatusOK)
}

// getUserInfo retrieves the authenticated user's information
// @Summary Get User Info
// @Description Retrieve the authenticated user's information
// @Tags AUTH
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SWSuccessResponse
// @Failure 401 {object} response.SWErrorResponse
// @Failure 500 {object} response.SWErrorResponse
// @Router /auth/userinfo [get]
func getUserInfo(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*jwt.Claims)
	var user models.User
	err := database.DB.Where("id = ?", userClaims.UserID).First(&user).Error
	if err != nil {
		return response.Fail(c, "USER_NOT_FOUND", "User not found", fiber.StatusNotFound)
	}

	token, err := jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to generate token", fiber.StatusInternalServerError)
	}

	now := time.Now()
	if err := database.DB.Model(&user).Update("last_login_at", &now).Error; err != nil {
    return response.Fail(c, "DB_ERROR", "Failed to update last login", fiber.StatusInternalServerError)
 }



	userInfo := schema.UserInfo{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		DisplayName:     user.DisplayName,
		Avatar:          user.Avatar,
		IsEmailVerified: user.IsEmailVerified,
		LastLoginAt:     user.LastLoginAt,
}
	return response.OK(c, fiber.Map{
		"user":  userInfo,
		"token": token,
	}, fiber.StatusOK)
}