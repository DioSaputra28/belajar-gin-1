package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	_ = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbName == "" {
		fmt.Println("ERROR: Missing required database environment variables")
		fmt.Printf("DB_HOST=%s, DB_PORT=%s, DB_NAME=%s\n", dbHost, dbPort, dbName)
		return nil
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("ERROR: Failed to connect to database: %v\n", err)
		fmt.Printf("Connection string: %s:***@tcp(%s:%s)/%s\n", dbUser, dbHost, dbPort, dbName)
		return nil
	}

	fmt.Println("Successfully connected to database")
	return db
}
