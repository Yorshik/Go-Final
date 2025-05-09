package server_test

import (
	"bytes"
	"encoding/json"
	"github.com/Yorshik/Go-Final/internal/database"
	"github.com/Yorshik/Go-Final/internal/models"
	"github.com/Yorshik/Go-Final/internal/server"
	"github.com/Yorshik/Go-Final/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestEnv() {
	os.Setenv("JWT_SECRET", "test-secret") // Устанавливаем переменные вручную
}

func setupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	_ = db.AutoMigrate(&models.User{}, &models.Expression{})
	database.DB = db
}

func TestRegisterHandler(t *testing.T) {
	setupTestEnv()
	setupTestDB()

	e := echo.New()
	reqBody := `{"username": "testuser", "password": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := server.RegisterHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "Registration successful", response["message"])
}

func TestLoginHandler(t *testing.T) {
	setupTestEnv()
	setupTestDB()

	// Регистрируем пользователя вручную
	hashed, _ := utils.HashPassword("testpass")
	user := models.User{Username: "testuser", Password: hashed}
	database.DB.Create(&user)

	e := echo.New()
	reqBody := `{"username": "testuser", "password": "testpass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := server.LoginHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NotEmpty(t, response["token"])
}
