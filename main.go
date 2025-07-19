package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/phonsing-Hub/GoLang/docs"
	"github.com/phonsing-Hub/GoLang/internal/config"
	"github.com/phonsing-Hub/GoLang/internal/database"
	"github.com/phonsing-Hub/GoLang/internal/middleware"
	"github.com/phonsing-Hub/GoLang/internal/routes"
	_ "github.com/phonsing-Hub/GoLang/internal/routes/api"
)

// @title GoLang API
// @version 1.0
// @description This is a sample API server using Fiber framework
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	config.LoadEnv()
	if err := database.Init(config.Env.DBUrl); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	app := fiber.New()
	app.Static("/static", "./static/uploads")

	config.SetupSwagger(app)

	app_v1 := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	app_v1.Use(cors.New(cors.Config{
		AllowOrigins: config.Env.CORSAllowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))
	routes.SetupRoutes(app_v1)
	routes.SetupMonitorRoute(app)

	app.Mount("/api/v1", app_v1)
	app.Listen(fmt.Sprintf(":%s", config.Env.AppPort))
}
