package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sukhantharot/go-service/handlers"
	"github.com/sukhantharot/go-service/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Add logging middleware
	router.Use(middleware.LoggingMiddleware())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})
	// Health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		// Check database connection
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "Database connection error",
				"error":   err.Error(),
			})
			return
		}

		// Ping database
		err = sqlDB.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "Database ping failed",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"database": "connected",
		})
	})

	// Public routes
	router.POST("/api/auth/register", handlers.Register)
	router.POST("/api/auth/login", handlers.Login)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		// User routes
		protected.GET("/users/me", handlers.GetCurrentUser)

		// Admin routes (example of role-based access)
		admin := protected.Group("/admin")
		admin.Use(middleware.RequirePermission("admin"))
		{
			admin.GET("/users", handlers.GetAllUsers)
			admin.POST("/roles", handlers.CreateRole)
			admin.POST("/permissions", handlers.CreatePermission)
		}
	}
}
