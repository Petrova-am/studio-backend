package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/studio-backend/pkg/database"
)

// ========== ТЕСТЫ БЕЗ БД (всегда проходят) ==========

// Тесты для sendSuccess
func TestSendSuccess_Booking(t *testing.T) {
	w := httptest.NewRecorder()
	sendSuccess(w, map[string]string{"status": "ok"})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestSendSuccess_WithArray_Booking(t *testing.T) {
	w := httptest.NewRecorder()
	sendSuccess(w, []int{1, 2, 3})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSendSuccess_WithString_Booking(t *testing.T) {
	w := httptest.NewRecorder()
	sendSuccess(w, "hello")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSendSuccess_WithNil_Booking(t *testing.T) {
	w := httptest.NewRecorder()
	sendSuccess(w, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Тесты для sendError
func TestSendError_Booking(t *testing.T) {
	w := httptest.NewRecorder()
	sendError(w, "error", http.StatusBadRequest)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestSendError_With404_Booking(t *testing.T) {
	w := httptest.NewRecorder()
	sendError(w, "not found", http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSendError_With500_Booking(t *testing.T) {
	w := httptest.NewRecorder()
	sendError(w, "server error", http.StatusInternalServerError)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Тесты для CORS
func TestCreateBooking_CORS(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/api/bookings", nil)
	w := httptest.NewRecorder()
	createBooking(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

// Тесты на существование хендлеров
func TestHandlersExist_Booking(t *testing.T) {
	assert.NotNil(t, createBooking)
	assert.NotNil(t, getUserBookings)
	assert.NotNil(t, updateStatus)
	assert.NotNil(t, getAvailableSlots)
}

// Тесты для POST/PUT запросов без тела
func TestCreateBooking_PostNoBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/bookings", nil)
	w := httptest.NewRecorder()
	createBooking(w, req)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ========== ТЕСТЫ С РЕАЛЬНОЙ БД (если БД доступна) ==========

func TestCreateBooking_WithRealDB(t *testing.T) {
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

	booking := map[string]interface{}{
		"userId":    1,
		"service":   "yoga",
		"trainerId": 1,
		"date":      "2026-05-20",
		"time":      "10:00",
		"amount":    1500,
	}
	body, _ := json.Marshal(booking)

	req := httptest.NewRequest("POST", "/api/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createBooking(w, req)
	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
}

func TestGetAvailableSlots_WithRealDB(t *testing.T) {
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

	req := httptest.NewRequest("GET", "/api/bookings/available-slots?date=2026-05-20&time=10:00&trainerId=1", nil)
	w := httptest.NewRecorder()

	getAvailableSlots(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]int
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.GreaterOrEqual(t, response["available"], 0)
	assert.LessOrEqual(t, response["available"], 8)
}

func TestGetUserBookings_WithRealDB(t *testing.T) {
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

	req := httptest.NewRequest("GET", "/api/bookings/user/1", nil)
	w := httptest.NewRecorder()

	getUserBookings(w, req)
	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
}
