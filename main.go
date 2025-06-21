package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/config"
	"github.com/phonsing-Hub/GoLang/database"
	"github.com/phonsing-Hub/GoLang/middleware"
	"github.com/phonsing-Hub/GoLang/routes"
	"log"
)

// @title           My API
// @version         0.1.1
// @description     This is a sample API GoFiber application.
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description Provide your API key in the 'X-API-Key' header.

func main() {
	// crate instance of config
	config.LoadEnv()

	if err := database.Init(config.Env.DBUrl); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	config.SetupSwagger(app)
	routes.SetupMonitorRoute(app)

	app.Use(middleware.FiberAccessLogger())

	routes.SetupRoutes(app)

	// start the server
	app.Listen(fmt.Sprintf(":%s", config.Env.AppPort))
}
