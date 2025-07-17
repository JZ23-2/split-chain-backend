package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("No .env file found or failed to load, relying on system env")
	}

	setupTestDB()

	exitCode := m.Run()

	cleanupTestDB()

	os.Exit(exitCode)

	os.Exit(m.Run())
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setupTestDB() {
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPass := getEnvOrDefault("DB_PASS", "")
	dbHost := getEnvOrDefault("DB_HOST", "127.0.0.1:3306")
	dbName := getEnvOrDefault("DB_NAME", "")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbName)

	fmt.Println(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	// Handle the error from AutoMigrate()
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Set the global database instance
	database.DB = db
}

func cleanupTestDB() {
	if database.DB != nil {
		// Clean up test data
		database.DB.Exec("DELETE FROM users")

		// Close connection
		sqlDB, err := database.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
		database.DB = nil
	}
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", controllers.RegisterUser)
	return r
}

func TestRegisterUser_Success(t *testing.T) {
	setupTestDB()
	defer cleanupTestDB()
	router := setupRouter()

	user := models.User{
		WalletAddress: "0.0.5161249",
	}

	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
}

func TestRegisterUser_DuplicateWallet(t *testing.T) {
	setupTestDB()
	defer cleanupTestDB()
	router := setupRouter()

	user := models.User{
		WalletAddress: "0.0.5161249",
	}

	// First registration
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should succeed
	assert.Equal(t, http.StatusOK, w.Code)

	// Second registration with same wallet
	body2, _ := json.Marshal(user)
	req2, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should fail with conflict
	assert.Equal(t, http.StatusConflict, w2.Code)
}
