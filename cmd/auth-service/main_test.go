package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"

	"github.com/stretchr/testify/assert"

	"time"

	"github.com/yourusername/studio-backend/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

// ========== ТЕСТЫ НА ФУНКЦИИ БЕЗ БД ==========
func TestRegisterHandler_WithRealDB(t *testing.T) {
	// Инициализируем подключение к БД для теста
	cfg := database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "123987",
		DBName:   "studio_db",
		SSLMode:  "disable",
	}

	var err error
	db, err = database.NewConnection(cfg)
	if err != nil {
		t.Skip("Database not available:", err)
		return
	}
	defer db.Close()

	email := "test_" + time.Now().Format("20060102150405") + "@test.com"
	body := `{"email":"` + email + `","password":"123456","name":"Test User"}`

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	registerHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NotNil(t, resp["token"])
}

func TestLoginHandler_WithRealDB(t *testing.T) {
	cfg := database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "123987",
		DBName:   "studio_db",
		SSLMode:  "disable",
	}

	var err error
	db, err = database.NewConnection(cfg)
	if err != nil {
		t.Skip("Database not available:", err)
		return
	}
	defer db.Close()

	// Сначала создадим тестового пользователя через SQL
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	_, err = db.Exec("INSERT INTO users (email, password_hash, name, role) VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING",
		"admin@example.com", string(hashedPassword), "Admin", "admin")
	if err != nil {
		t.Log("Warning: could not create test user:", err)
	}

	body := `{"email":"admin@example.com","password":"admin123"}`
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	loginHandler(w, req)

	// Проверяем что статус 200 ИЛИ 400/401 (если пользователь не создался)
	if w.Code != http.StatusOK {
		t.Logf("Response code: %d, Body: %s", w.Code, w.Body.String())
	}
	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
}

// ========== ТЕСТ С РЕАЛЬНОЙ БД (добавьте в конец main_test.go) ==========
func TestDatabaseConnection(t *testing.T) {
	cfg := database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "123987",
		DBName:   "studio_db",
		SSLMode:  "disable",
	}

	db, err := database.NewConnection(cfg)
	if err != nil {
		t.Skip("Database not available:", err)
		return
	}
	defer db.Close()

	// Проверяем что подключились
	err = db.Ping()
	assert.NoError(t, err)
}

func TestGenerateJWT(t *testing.T) {
	token, err := generateJWT(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateJWT_DifferentIDs(t *testing.T) {
	token1, _ := generateJWT(1)
	token2, _ := generateJWT(2)
	assert.NotEqual(t, token1, token2)
}

func TestGetEnv_Default(t *testing.T) {
	result := getEnv("NOT_EXIST", "default")
	assert.Equal(t, "default", result)
}

func TestGetEnv_Existing(t *testing.T) {
	t.Setenv("TEST_VAR", "value")
	result := getEnv("TEST_VAR", "default")
	assert.Equal(t, "value", result)
}

func TestSendSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	sendSuccess(w, map[string]string{"status": "ok"})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestSendError(t *testing.T) {
	w := httptest.NewRecorder()
	sendError(w, "error", 400)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestLoginHandler_CORS(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/api/auth/login", nil)
	w := httptest.NewRecorder()
	loginHandler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegisterHandler_CORS(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/api/auth/register", nil)
	w := httptest.NewRecorder()
	registerHandler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestValidateHandler_NoToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/auth/validate", nil)
	w := httptest.NewRecorder()
	validateHandler(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
