package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/sukhantharot/go-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func cleanConnectionString(connStr string) string {
	// If it's a URL format
	if strings.HasPrefix(connStr, "postgres://") || strings.HasPrefix(connStr, "postgresql://") {
		u, err := url.Parse(connStr)
		if err != nil {
			log.Printf("Warning: Could not parse DATABASE_URL as URL: %v", err)
			return connStr
		}
		q := u.Query()
		// Remove problematic parameters
		q.Del("schema")
		q.Del("connection_limit")
		u.RawQuery = q.Encode()
		return u.String()
	}

	// If it's a key=value format
	parts := strings.Split(connStr, " ")
	validParts := make([]string, 0)
	for _, part := range parts {
		// Skip problematic parameters
		if !strings.HasPrefix(part, "schema=") &&
			!strings.HasPrefix(part, "connection_limit=") {
			validParts = append(validParts, part)
		}
	}
	return strings.Join(validParts, " ")
}

func InitDB() *gorm.DB {
	log.Println("Starting database initialization...")
	log.Printf("Running in environment: %s", os.Getenv("APP_ENV"))
	log.Printf("RAILWAY_ENVIRONMENT: %s", os.Getenv("RAILWAY_ENVIRONMENT"))

	// Print all environment variables for debugging
	log.Println("Environment variables:")
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			key := pair[0]
			value := pair[1]
			// Mask sensitive information
			if strings.Contains(strings.ToLower(key), "password") ||
				strings.Contains(strings.ToLower(key), "secret") ||
				strings.Contains(strings.ToLower(key), "token") ||
				strings.Contains(strings.ToLower(key), "database_url") {
				value = "****"
			}
			log.Printf("%s=%s", key, value)
		}
	}

	var dsn string
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL != "" {
		log.Println("Using DATABASE_URL for connection")
		dsn = cleanConnectionString(dbURL)
		log.Println("Connection string cleaned and prepared")
	} else if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
		log.Fatal("Error: Running on Railway but DATABASE_URL is not set. Please check:\n" +
			"1. PostgreSQL database is provisioned in Railway\n" +
			"2. Database is linked to your service\n" +
			"3. DATABASE_URL variable is set in service variables")
	} else {
		log.Println("Using individual connection parameters")
		// Local development with individual parameters
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")

		// Check if any required parameter is missing
		missingParams := []string{}
		if host == "" {
			missingParams = append(missingParams, "DB_HOST")
		}
		if user == "" {
			missingParams = append(missingParams, "DB_USER")
		}
		if dbPassword == "" {
			missingParams = append(missingParams, "DB_PASSWORD")
		}
		if dbname == "" {
			missingParams = append(missingParams, "DB_NAME")
		}
		if port == "" {
			missingParams = append(missingParams, "DB_PORT")
		}

		if len(missingParams) > 0 {
			log.Fatalf("Missing required database parameters: %s", strings.Join(missingParams, ", "))
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, dbPassword, dbname, port)
	}

	log.Println("Attempting database connection...")

	// Configure GORM with detailed logging
	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
				Colorful: false,
			},
		),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		log.Fatal("Failed to connect to database. Please check your configuration.")
	}

	log.Println("Successfully connected to database")

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database ping successful")

	// Auto Migrate the schema
	log.Println("Starting database migration...")
	if err := db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}); err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
	log.Println("Database migration completed successfully")

	DB = db
	return db
}
