package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/oat/go-service/handlers"
	"github.com/oat/go-service/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
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