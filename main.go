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
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	routes.SetupRoutes(app)
	routes.SetupMonitorRoute(app)
	app.Listen(fmt.Sprintf(":%s", config.Env.AppPort))
}
