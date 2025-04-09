package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sukhantharot/go-service/config"
	"github.com/sukhantharot/go-service/routes"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database
	db := config.InitDB()

	// Create Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, db)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
} 