package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/database"
	"github.com/phonsing-Hub/GoLang/database/models"
	"github.com/phonsing-Hub/GoLang/utils/helper"
)



func SetupUserRoutes(router fiber.Router) {
	userGroup := router.Group("/users")
	userGroup.Get("/", get_users)
	userGroup.Get("/:id", get_user_id)
}

// @Summary      Get all users
// @Description  Retrieve all users from the database
// @Tags         Users
// @Produce      json
// @Success      200 {object} response.SWListSuccessResponse
// @Failure      500 {object} response.SWErrorResponse
// @Router       /users [get]
// @Security     BearerAuth
func get_users(c *fiber.Ctx) error {
	return helper.FindAll[models.Users](c, database.DB)
}

// @Summary      Get user by ID
// @Description  Retrieve all users from the database
// @Tags         Users
// @Produce      json
// @Success      200 {object} response.SWSuccessResponse
// @Failure      500 {object} response.SWErrorResponse
// @Router       /users/{id} [get]
// @Param        id path int true "User ID"
// @Security     BearerAuth
func get_user_id(c *fiber.Ctx) error {
	return helper.FindByID[models.Users](c, database.DB)
} 