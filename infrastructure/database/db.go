package database

import (
	"fmt"
	"log"
	"wms-be/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB: inisialisasi koneksi dan migrasi
func InitDB() {
	if DB != nil {
		return
	}

	dbURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.GetEnv("DB_HOST"),
		config.GetEnv("DB_PORT"),
		config.GetEnv("DB_USER"),
		config.GetEnv("DB_PASSWORD"),
		config.GetEnv("DB_NAME"),
		config.GetEnv("DB_SSL_MODE"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
}

// GetDB: kembalikan pointer *gorm.DB, panic jika belum di-init
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("database not initialized! call database.InitDB() first")
	}
	return DB
}
