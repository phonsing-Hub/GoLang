package database

import (
  "gorm.io/gorm"
  "gorm.io/driver/postgres"

)

var DB *gorm.DB

func Init(dsn string) error {

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(0) 

	return nil
}