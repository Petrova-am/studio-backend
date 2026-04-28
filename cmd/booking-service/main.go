package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/studio-backend/pkg/database"
	"github.com/yourusername/studio-backend/pkg/logger"
)

var db *sql.DB

type Booking struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	Service   string `json:"service"`
	TrainerID int    `json:"trainerId"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Status    string `json:"status"`
	Amount    int    `json:"amount"`
}

func main() {
	logger.InitLogger()

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
		logger.Error("Failed to connect to database", err)
		return
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/bookings", createBooking).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/bookings/user/{userId}", getUserBookings).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/bookings/{id}/status", updateStatus).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/bookings/available-slots", getAvailableSlots).Methods("GET", "OPTIONS")

	logger.Info("Booking service running on port 8082")
	http.ListenAndServe(":8082", r)
}

func createBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	var booking Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		sendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO bookings (user_id, service, trainer_id, date, time, status, amount) 
              VALUES ($1, $2, $3, $4, $5, 'pending', $6) RETURNING id`
	err := db.QueryRow(query, booking.UserID, booking.Service, booking.TrainerID,
		booking.Date, booking.Time, booking.Amount).Scan(&booking.ID)

	if err != nil {
		logger.Error("Failed to create booking", err)
		sendError(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	sendSuccess(w, booking)
}

func getUserBookings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userId"])

	rows, err := db.Query("SELECT id, service, trainer_id, date, time, status FROM bookings WHERE user_id = $1", userID)
	if err != nil {
		sendError(w, "Failed to get bookings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		rows.Scan(&b.ID, &b.Service, &b.TrainerID, &b.Date, &b.Time, &b.Status)
		bookings = append(bookings, b)
	}

	sendSuccess(w, bookings)
}

func updateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var req struct {
		Status string `json:"status"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	_, err := db.Exec("UPDATE bookings SET status = $1 WHERE id = $2", req.Status, id)
	if err != nil {
		sendError(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	sendSuccess(w, map[string]string{"status": "updated"})
}

func getAvailableSlots(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	time := r.URL.Query().Get("time")
	trainerID := r.URL.Query().Get("trainerId")

	var count int
	query := "SELECT COUNT(*) FROM bookings WHERE date = $1 AND time = $2 AND trainer_id = $3"
	db.QueryRow(query, date, time, trainerID).Scan(&count)

	available := 8 - count
	sendSuccess(w, map[string]int{"available": available})
}

func sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
