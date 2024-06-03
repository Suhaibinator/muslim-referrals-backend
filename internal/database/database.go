package database

import (
	"log"

	"muslim-referrals-backend/configs"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    var err error
    DB, err = gorm.Open(sqlite.Open(configs.DatabasePath), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    log.Println("Database connection successful.")
}

func CloseDatabase() {
	db, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to close database:", err)
	}
	db.Close()
	log.Println("Database connection closed.")
}

