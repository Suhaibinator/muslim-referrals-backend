package database

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    var err error
    DB, err = gorm.Open(sqlite.Open(os.Getenv("DATABASE")), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    log.Println("Database connection successful.")
}