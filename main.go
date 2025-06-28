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

func main() {
	config.LoadEnv()
	if err := database.Init(config.Env.DBUrl); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
  app := fiber.New()
	app_v1 := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	routes.SetupRoutes(app_v1)
	routes.SetupMonitorRoute(app)
	
	app.Mount("/api/v1", app_v1)
	app.Listen(fmt.Sprintf(":%s", config.Env.AppPort))
}
