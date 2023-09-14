package database

import (
	"fmt"

	"question-and-answers/pkg/config"
	"question-and-answers/pkg/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenConnection(appConfig config.AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		appConfig.DatabaseHost,
		appConfig.DatabaseUser,
		appConfig.DatabasePassword,
		appConfig.DatabaseName,
		appConfig.DatabasePort,
		appConfig.DatabaseTimezone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Question{}); err != nil {
		return nil, err
	}

	return db, nil
}
