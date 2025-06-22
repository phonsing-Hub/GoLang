package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/database"
	"github.com/phonsing-Hub/GoLang/database/models"
	"github.com/phonsing-Hub/GoLang/database/schema"
	"github.com/phonsing-Hub/GoLang/middleware"
	"github.com/phonsing-Hub/GoLang/pkg/auth"
	"github.com/phonsing-Hub/GoLang/utils/helper"
	"github.com/phonsing-Hub/GoLang/utils/response"
)

func SetupUserRoutes(router fiber.Router) {
	userGroup := router.Group("/users")
	userGroup.Use(middleware.JWTAuthMiddleware())

	userGroup.Get("/", get_users)
	userGroup.Get("/:id", get_user_id)
	userGroup.Post("/", create_user)
	userGroup.Put("/:id", update_user)
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