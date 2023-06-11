package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(DatabaseURL string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(DatabaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// db.AutoMigrate(&models.Employee{}, &models.Department{}, &models.Position{})

	return db
}
