package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/phonsing-Hub/GoLang/internal/middleware"
	"github.com/phonsing-Hub/GoLang/internal/routes/api"

)

func SetupRoutes(app *fiber.App) {
	app.Use(middleware.FiberAccessLogger())
	app.Use(middleware.ZapLogger())

	api.SetupAuthRoutes(app)
	api.SetupUserhRoutes(app)

}


func SetupMonitorRoute(app *fiber.App) {
	app.Get("/monitoring", monitor.New())
}
