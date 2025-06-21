package main

import (
	"fmt"
	"github.com/phonsing-Hub/GoLang/config"
	"github.com/phonsing-Hub/GoLang/database"
	"github.com/phonsing-Hub/GoLang/database/models"
	"github.com/phonsing-Hub/GoLang/database/views"
	"log"
)

func main() {
	config.LoadEnv()
	if err := database.Init(config.Env.DBUrl); err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}

	if err := database.DB.AutoMigrate(models.All()...); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

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

	log.Println("Migration completed successfully")
}
