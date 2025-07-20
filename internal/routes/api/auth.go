package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/internal/config"
	"github.com/phonsing-Hub/GoLang/internal/database"
	"github.com/phonsing-Hub/GoLang/internal/database/models"
	"github.com/phonsing-Hub/GoLang/internal/database/schema"
	"github.com/phonsing-Hub/GoLang/internal/middleware"
	"github.com/phonsing-Hub/GoLang/internal/utils/helper"
	"github.com/phonsing-Hub/GoLang/internal/utils/response"
	"github.com/phonsing-Hub/GoLang/pkg/auth"
	"github.com/phonsing-Hub/GoLang/pkg/jwt"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

func SetupAuthRoutes(router *fiber.App) {
	authGroup := router.Group("/auth")
	authGroup.Post("/register", registerUser)
	authGroup.Post("/login", loginUser)
	// Google OAuth routes
	authGroup.Post("/google", googleLogin)
	authGroup.Post("/google/callback", googleCallback)

	authGroup.Get("/userinfo", middleware.JWTAuthMiddleware(), getUserInfo)
}

func exchangeCodeForToken(code string) (*schema.GoogleTokenResponse, error) {
	values := url.Values{}
	values.Set("code", code)
	values.Set("client_id", config.Env.GoogleClientID)
	values.Set("client_secret", config.Env.GoogleClientSecret)
	values.Set("redirect_uri", config.Env.GoogleRedirectURI)
	values.Set("grant_type", config.Env.GoogleGrantType)

	resp, err := http.PostForm(config.Env.GoogleTokenURL, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenRes schema.GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return nil, err
	}
	return &tokenRes, nil

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

// googleLogin redirects user to Google OAuth
// @Summary Google OAuth Login
// @Description Redirect to Google OAuth login
// @Tags AUTH
// @Accept json
// @Produce json
// @Success 200 {object} response.SWSuccessResponse
// @Router /auth/google [post]
func googleLogin(c *fiber.Ctx) error {
	authProvider := "google"
	authType := "oauth"
	now := time.Now()
	type Req struct {
		Credential string `json:"credential"`
	}
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, "BODY_PARSE_ERROR", "Failed to parse request body", fiber.StatusBadRequest)
	}

	payload, err := idtoken.Validate(context.Background(), req.Credential, config.Env.GoogleClientID)
	if err != nil {
		return response.Fail(c, "INVALID_TOKEN", "Invalid token", fiber.StatusUnauthorized)
	}

	googleID := payload.Claims["sub"].(string)
	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	FirstName := payload.Claims["given_name"].(string)
	LastName := payload.Claims["family_name"].(string)
	picture := payload.Claims["picture"].(string)

	// Check if user already exists
	var existingAuth models.UserAuthMethod
	result := database.DB.Where("auth_provider = ? AND provider_id = ?", authProvider, googleID).First(&existingAuth)

	if result.Error == nil {
		var user models.User
		if err := database.DB.Where("id = ?", existingAuth.UserID).First(&user).Error; err != nil {
			return response.Fail(c, "USER_NOT_FOUND", "User not found", fiber.StatusNotFound)
		}
		if err := database.DB.Model(&user).Update("last_login_at", &now).Error; err != nil {
			return response.Fail(c, "DB_ERROR", "Failed to update last login", fiber.StatusInternalServerError)
		}
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
		}, fiber.StatusOK)
	}

	newUser := models.User{
		Email:           email,
		DisplayName:     name,
		Avatar:          picture,
		FirstName:       FirstName,
		LastName:        LastName,
		LastLoginAt:     &now,
		IsEmailVerified: true,
	}
	if err := database.DB.Create(&newUser).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to create user", fiber.StatusInternalServerError)
	}
	userAuth := models.UserAuthMethod{
		UserID:       newUser.ID,
		AuthType:     authType,
		AuthProvider: &authProvider,
		ProviderID:   &googleID,
		IsPrimary:    true,
	}
	if err := database.DB.Create(&userAuth).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to create user auth method", fiber.StatusInternalServerError)
	}
	token, err := jwt.GenerateToken(newUser.ID, newUser.Email)
	if err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to generate token", fiber.StatusInternalServerError)
	}

	userInfo := schema.UserInfo{
		ID:              newUser.ID,
		Email:           newUser.Email,
		FirstName:       newUser.FirstName,
		LastName:        newUser.LastName,
		DisplayName:     newUser.DisplayName,
		Avatar:          newUser.Avatar,
		IsEmailVerified: newUser.IsEmailVerified,
		LastLoginAt:     newUser.LastLoginAt,
	}

	return response.OK(c, fiber.Map{
		"user":  userInfo,
		"token": token,
	}, fiber.StatusCreated)
}

// googleCallback handles Google OAuth callback
// @Summary Google OAuth Callback
// @Description Handle Google OAuth callback and exchange code for tokens
// @Tags AUTH
// @Accept json
// @Produce json
// @Param request body object{code=string} true "Authorization Code Request" example({"code": "4/0AX4XfWh..."})
// @Success 200 {object} response.SWSuccessResponse
// @Failure 400 {object} response.SWErrorResponse
// @Failure 401 {object} response.SWErrorResponse
// @Failure 500 {object} response.SWErrorResponse
// @Router /auth/google/callback [post]
func googleCallback(c *fiber.Ctx) error {
	authProvider := "google"
	authType := "oauth"
	now := time.Now()
	type Req struct {
		Code string `json:"code"`
	}
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, "BODY_PARSE_ERROR", "Failed to parse request body", fiber.StatusBadRequest)
	}

	// Step 1: แลก code เป็น token
	tokenResp, err := exchangeCodeForToken(req.Code)
	if err != nil {
		return response.Fail(c, "TOKEN_EXCHANGE_ERROR", "Failed to exchange code", fiber.StatusUnauthorized)
	}

	// Step 2: Decode id_token (JWT)
	fmt.Println("ID Token:", tokenResp.IDToken)
	payload, err := idtoken.Validate(context.Background(), tokenResp.IDToken, config.Env.GoogleClientID)
	if err != nil {
		return response.Fail(c, "INVALID_TOKEN", "Invalid ID Token", fiber.StatusUnauthorized)
	}

	googleID := payload.Claims["sub"].(string)
	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	FirstName := payload.Claims["given_name"].(string)
	LastName := payload.Claims["family_name"].(string)
	picture := payload.Claims["picture"].(string)

	// Check if user already exists
	var existingAuth models.UserAuthMethod
	result := database.DB.Where("auth_provider = ? AND provider_id = ?", authProvider, googleID).First(&existingAuth)

	if result.Error == nil {
		var user models.User
		if err := database.DB.Where("id = ?", existingAuth.UserID).First(&user).Error; err != nil {
			return response.Fail(c, "USER_NOT_FOUND", "User not found", fiber.StatusNotFound)
		}
		if err := database.DB.Model(&user).Update("last_login_at", &now).Error; err != nil {
			return response.Fail(c, "DB_ERROR", "Failed to update last login", fiber.StatusInternalServerError)
		}
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
		}, fiber.StatusOK)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var existingUser models.User
		if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
			if emailVerified, ok := payload.Claims["email_verified"].(bool); ok && emailVerified {

				userAuth := models.UserAuthMethod{
					UserID:       existingUser.ID,
					AuthType:     authType,
					AuthProvider: &authProvider,
					ProviderID:   &googleID,
					IsPrimary:    false,
				}

				if err := database.DB.Create(&userAuth).Error; err != nil {
					return response.Fail(c, "DB_ERROR", "Failed to link Google account", fiber.StatusInternalServerError)
				}

				_ = database.DB.Model(&existingUser).Updates(map[string]interface{}{
					"last_login_at":     now,
					"is_email_verified": true,
				})
				
				token, _ := jwt.GenerateToken(existingUser.ID, existingUser.Email)
				userInfo := schema.UserInfo{
					ID:              existingUser.ID,
					Email:           existingUser.Email,
					FirstName:       existingUser.FirstName,
					LastName:        existingUser.LastName,
					DisplayName:     existingUser.DisplayName,
					Avatar:          existingUser.Avatar,
					IsEmailVerified: existingUser.IsEmailVerified,
					LastLoginAt:     existingUser.LastLoginAt,
				}
				return response.OK(c, fiber.Map{
					"user":  userInfo,
					"token": token,
				}, fiber.StatusOK)
			}
		}
	}

	newUser := models.User{
		Email:           email,
		DisplayName:     name,
		Avatar:          picture,
		FirstName:       FirstName,
		LastName:        LastName,
		LastLoginAt:     &now,
		IsEmailVerified: true,
	}
	if err := database.DB.Create(&newUser).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to create user", fiber.StatusInternalServerError)
	}
	userAuth := models.UserAuthMethod{
		UserID:       newUser.ID,
		AuthType:     authType,
		AuthProvider: &authProvider,
		ProviderID:   &googleID,
		IsPrimary:    true,
	}
	if err := database.DB.Create(&userAuth).Error; err != nil {
		return response.Fail(c, "DATABASE_ERROR", "Failed to create user auth method", fiber.StatusInternalServerError)
	}
	token, err := jwt.GenerateToken(newUser.ID, newUser.Email)
	if err != nil {
		return response.Fail(c, "INTERNAL_ERROR", "Failed to generate token", fiber.StatusInternalServerError)
	}

	userInfo := schema.UserInfo{
		ID:              newUser.ID,
		Email:           newUser.Email,
		FirstName:       newUser.FirstName,
		LastName:        newUser.LastName,
		DisplayName:     newUser.DisplayName,
		Avatar:          newUser.Avatar,
		IsEmailVerified: newUser.IsEmailVerified,
		LastLoginAt:     newUser.LastLoginAt,
	}

	return response.OK(c, fiber.Map{
		"user":  userInfo,
		"token": token,
	}, fiber.StatusCreated)
}
