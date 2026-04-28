package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourusername/studio-backend/internal/models"
	"github.com/yourusername/studio-backend/pkg/database"
	"github.com/yourusername/studio-backend/pkg/logger"
)

var db *sql.DB

func main() {
	logger.InitLogger()

	cfg := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "123987"),
		DBName:   getEnv("DB_NAME", "studio_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	var err error
	db, err = database.NewConnection(cfg)
	if err != nil {
		logger.Error("Failed to connect to database", err)
		os.Exit(1)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/auth/login", loginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/register", registerHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/validate", validateHandler).Methods("GET")

	logger.Infof("Auth service running on port %s", getEnv("PORT", "8081"))
	http.ListenAndServe(":"+getEnv("PORT", "8081"), r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user struct {
		ID           int
		Email        string
		PasswordHash string
		Name         string
		Role         string
	}

	query := "SELECT id, email, password_hash, name, role FROM users WHERE email = $1"
	err := db.QueryRow(query, req.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		logger.Error("Database error", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := generateJWT(user.ID)
	sendSuccess(w, models.Response{
		Token: token,
		User: models.User{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var userID int
	query := "INSERT INTO users (email, password_hash, name, role) VALUES ($1, $2, $3, 'user') RETURNING id"
	err = db.QueryRow(query, req.Email, string(hashedPassword), req.Name).Scan(&userID)

	if err != nil {
		logger.Error("Failed to create user", err)
		sendError(w, "Email already exists", http.StatusConflict)
		return
	}

	token, _ := generateJWT(userID)
	sendSuccess(w, models.Response{
		Token: token,
		User: models.User{
			ID:    userID,
			Email: req.Email,
			Name:  req.Name,
			Role:  "user",
		},
	})
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		sendError(w, "No token provided", http.StatusUnauthorized)
		return
	}

	sendSuccess(w, map[string]bool{"valid": true})
}

func generateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
	})
	return token.SignedString([]byte("your-secret-key"))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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
