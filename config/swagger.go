package config

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func SetupSwagger(app *fiber.App) {
	swaggerConfig := swagger.Config{
		BasePath: "/api/v1",
		FilePath: "./docs/swagger.json",
		Path:     "docs",
		Title:    "Fiber API Docs",
	}
	app.Use(swagger.New(swaggerConfig))
}
