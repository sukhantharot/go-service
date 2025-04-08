package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestResponse represents a generic test response
type TestResponse struct {
	Status  int         `json:"status,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// MakeRequest is a helper function to make HTTP requests in tests
func MakeRequest(t *testing.T, router *gin.Engine, method, url string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		assert.NoError(t, err)
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// AssertResponse is a helper function to assert HTTP responses
func AssertResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedError string) {
	assert.Equal(t, expectedStatus, w.Code)

	var response TestResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	if expectedError != "" {
		assert.Equal(t, expectedError, response.Error)
	}
}

// CreateTestUser is a helper function to create a test user
func CreateTestUser(t *testing.T, db *gorm.DB, email, password string) *models.User {
	user := &models.User{
		Email:     email,
		Password:  password,
		FirstName: "Test",
		LastName:  "User",
		RoleID:    1,
	}

	err := db.Create(user).Error
	assert.NoError(t, err)

	return user
}

// GenerateTestToken is a helper function to generate a JWT token for testing
func GenerateTestToken(t *testing.T, userID uint, roleID uint) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role_id": roleID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	assert.NoError(t, err)

	return tokenString
} 