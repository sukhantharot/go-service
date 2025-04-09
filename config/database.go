package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sukhantharot/go-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	var dsn string

	// Check if DATABASE_URL is set (Railway style)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
		log.Println("Using DATABASE_URL for connection")
	} else if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
		// If we're on Railway but no DATABASE_URL, this is an error
		log.Fatal("Error: Running on Railway but DATABASE_URL is not set")
	} else {
		// Local development with individual parameters
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")

		// Check if any required parameter is missing
		if host == "" || user == "" || password == "" || dbname == "" || port == "" {
			log.Printf("Missing database configuration: host=%s, user=%s, dbname=%s, port=%s",
				host, user, dbname, port)
			log.Fatal("Please set all required database environment variables or DATABASE_URL")
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, password, dbname, port)
		log.Println("Using individual parameters for connection")
	}

	// Log the connection string with password masked
	maskedDSN := strings.Replace(dsn, os.Getenv("DB_PASSWORD"), "****", 1)
	log.Printf("Connecting to database with: %s", maskedDSN)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Successfully connected to database")

	// Auto Migrate the schema
	if err := db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}); err != nil {
		log.Fatal("Failed to auto-migrate database schema: ", err)
	}

	DB = db
	return db
}
