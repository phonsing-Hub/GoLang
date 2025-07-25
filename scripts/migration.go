package main

import (
	"fmt"
	"log"

	"github.com/phonsing-Hub/GoLang/internal/config"
	"github.com/phonsing-Hub/GoLang/internal/database"
	"github.com/phonsing-Hub/GoLang/internal/database/models"
	"github.com/phonsing-Hub/GoLang/internal/database/views"
)

func main() {
	config.LoadEnv()
	if err := database.Init(config.Env.DBUrl); err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}

	modelsList := models.All()
	log.Printf("Starting migration for %d models:", len(modelsList))
	for i, model := range modelsList {
		modelName := fmt.Sprintf("%T", model)
		log.Printf("[%d/%d] Migrating: %s", i+1, len(modelsList), modelName)
	}

	if err := database.DB.AutoMigrate(models.All()...); err != nil {
		log.Fatalf("migration failed: %v", err)
	}


	log.Println("Managing database views...")
	for name := range views.Views {
		dropSQL := fmt.Sprintf("DROP VIEW IF EXISTS %s CASCADE;", name)
		if err := database.DB.Exec(dropSQL).Error; err != nil {
			log.Fatalf("failed to drop view %s: %v", name, err)
		}
		log.Printf("View dropped: %s", name)
	}

	for name, query := range views.Views {
		if err := database.DB.Exec(query).Error; err != nil {
			log.Fatalf("failed to create view %s: %v", name, err)
		}
		log.Printf("View created: %s", name)
	}

	//default userstatus
	defaultStatuses := []*models.UserStatus{
		{Name: "active", DisplayName: "Active", Description: "User is active", Color: "#28a745", IsActive: true, Position: 1},
		{Name: "inactive", DisplayName: "Inactive", Description: "User is inactive", Color: "#6c757d", IsActive: true, Position: 2},
		{Name: "suspended", DisplayName: "Suspended", Description: "User is suspended", Color: "#dc3545", IsActive: true, Position: 3},
		{Name: "pending_verification", DisplayName: "Pending Verification", Description: "User is pending verification", Color: "#ffc107", IsActive: true, Position: 4},
	}
	if err := database.DB.Create(defaultStatuses).Error; err != nil {
		log.Fatalf("failed to create default user statuses: %v", err)
	}

	log.Println("Migration completed successfully")
}


