package database

import (
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbDriver struct {
	mu sync.RWMutex
	db *gorm.DB
}

func NewDbDriver(dbPath string) *DbDriver {

	gormDb, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connection successful.")
	return &DbDriver{db: gormDb}
}

func (dbd *DbDriver) CloseDatabase() {
	sqlDb, err := dbd.db.DB()
	if err != nil {
		log.Fatal("Failed to close database:", err)
	}
	sqlDb.Close()
	log.Println("Database connection closed.")
}

func (db *DbDriver) AddRecord(record interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db.Create(record)
}

func (db *DbDriver) GetUser(userId uint64) *User {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var user User
	db.db.First(&user, userId)
	return &user
}
