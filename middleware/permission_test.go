package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oat/go-service/config"
	"github.com/oat/go-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Setup test database
	db, err := gorm.Open(postgres.Open("host=localhost user=test password=test dbname=test_db port=5432 sslmode=disable"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate schema
	err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})
	require.NoError(t, err)

	return db
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	config.DB = db
	return router
}

func TestRequirePermission(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := setupTestRouter(db)

	// Create test role and permission
	permission := models.Permission{
		Name:        "admin",
		Description: "Admin permission",
	}
	db.Create(&permission)

	role := models.Role{
		Name:        "admin",
		Description: "Admin role",
		Permissions: []models.Permission{permission},
	}
	db.Create(&role)

	// Test cases
	tests := []struct {
		name           string
		setupContext   func(c *gin.Context)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "user with required permission",
			setupContext: func(c *gin.Context) {
				c.Set("role_id", role.ID)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "user without role_id in context",
			setupContext: func(c *gin.Context) {
				// Do nothing, role_id not set
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Role ID not found in context",
		},
		{
			name: "user with non-existent role",
			setupContext: func(c *gin.Context) {
				c.Set("role_id", 999) // Non-existent role ID
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Role not found",
		},
		{
			name: "user without required permission",
			setupContext: func(c *gin.Context) {
				// Create role without admin permission
				noPermissionRole := models.Role{
					Name:        "user",
					Description: "Regular user role",
				}
				db.Create(&noPermissionRole)
				c.Set("role_id", noPermissionRole.ID)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "Permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup route with middleware
			router.GET("/test", func(c *gin.Context) {
				tt.setupContext(c)
			}, RequirePermission("admin"), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			}
		})
	}
} 