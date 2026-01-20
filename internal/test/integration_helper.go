package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/DioSaputra28/belajar-gin-1/internal/addresses"
	"github.com/DioSaputra28/belajar-gin-1/internal/contacts"
	"github.com/DioSaputra28/belajar-gin-1/internal/users"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SetupTestDB creates a database connection for integration tests
func SetupTestDB(t *testing.T) *gorm.DB {
	// Try to load .env from different possible locations
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load(".env")

	// Create database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect database: %v\nDSN: %s:%s@tcp(%s:%s)/%s",
			err,
			os.Getenv("DB_USER"),
			"***",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}

	// Drop existing tables to ensure clean migration
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("DROP TABLE IF EXISTS addresses")
	db.Exec("DROP TABLE IF EXISTS contacts")
	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// Auto migrate all tables - ORDER MATTERS (foreign keys)
	// Users must be created first before contacts (contacts references users)
	err = db.AutoMigrate(&users.User{})
	if err != nil {
		t.Fatalf("Failed to migrate users table: %v", err)
	}

	err = db.AutoMigrate(&contacts.Contact{})
	if err != nil {
		t.Fatalf("Failed to migrate contacts table: %v", err)
	}

	err = db.AutoMigrate(&addresses.Address{})
	if err != nil {
		t.Fatalf("Failed to migrate addresses table: %v", err)
	}

	return db
}

// CleanupTestDB removes all test data from tables
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	// Delete in correct order (foreign key constraints)
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("TRUNCATE TABLE addresses")
	db.Exec("TRUNCATE TABLE contacts")
	db.Exec("TRUNCATE TABLE users")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

// CleanupTestDBAfter is a helper to defer cleanup
func CleanupTestDBAfter(t *testing.T, db *gorm.DB) {
	CleanupTestDB(t, db)
}
