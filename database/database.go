package database

import (
	"ticket-system/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(databasePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}, &models.Ticket{}); err != nil {
		return nil, err
	}

	return db, nil
}
