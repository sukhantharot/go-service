package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestJWTAuth(t *testing.T) {
	// Set JWT secret
	os.Setenv("JWT_SECRET", "test_secret")

	// Test cases
	tests := []struct {
		name           string
		setupAuth      func() string
		expectedStatus int
		expectedError  string
	}{
		{
			name: "valid token",
			setupAuth: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id": 1,
					"role_id": 1,
					"exp":     time.Now().Add(time.Hour * 24).Unix(),
				})
				tokenString, _ := token.SignedString([]byte("test_secret"))
				return tokenString
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing token",
			setupAuth: func() string {
				return ""
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization header is required",
		},
		{
			name: "invalid token format",
			setupAuth: func() string {
				return "invalid_token_format"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization header format must be Bearer {token}",
		},
		{
			name: "expired token",
			setupAuth: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id": 1,
					"role_id": 1,
					"exp":     time.Now().Add(-time.Hour).Unix(), // Expired
				})
				tokenString, _ := token.SignedString([]byte("test_secret"))
				return tokenString
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupTestRouter()
			router.GET("/test", JWTAuth(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.setupAuth() != "" {
				req.Header.Set("Authorization", "Bearer "+tt.setupAuth())
			}
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