package api

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/phonsing-Hub/GoLang/internal/database"
	// "github.com/phonsing-Hub/GoLang/internal/database/models"
	// "github.com/phonsing-Hub/GoLang/internal/database/schema"
	// "github.com/phonsing-Hub/GoLang/internal/middleware"
	// "github.com/phonsing-Hub/GoLang/internal/utils/helper"
	// "github.com/phonsing-Hub/GoLang/internal/utils/response"
	// "github.com/phonsing-Hub/GoLang/pkg/auth"
	// "github.com/phonsing-Hub/GoLang/pkg/jwt"
	// "gorm.io/gorm"
)

func SetupUserhRoutes(router *fiber.App) {
	userGroup := router.Group("/user")
	// userGroup.Use(middleware.JWTAuthMiddleware())
	userGroup.Post("/", registerUser)
}

