package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sukhantharot/go-service/config"
	"github.com/sukhantharot/go-service/models"
)

func RequirePermission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role ID not found in context"})
			c.Abort()
			return
		}

		var role models.Role
		if err := config.DB.Preload("Permissions").First(&role, roleID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}

		hasPermission := false
		for _, permission := range role.Permissions {
			if permission.Name == permissionName {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
