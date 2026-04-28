package models

// User - модель пользователя
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

// Booking - модель бронирования
type Booking struct {
	ID            int    `json:"id"`
	UserID        int    `json:"userId"`
	Service       string `json:"service"`
	TrainerID     int    `json:"trainerId"`
	Date          string `json:"date"`
	Time          string `json:"time"`
	Status        string `json:"status"`
	PaymentStatus string `json:"paymentStatus"`
	Amount        int    `json:"amount"`
	Comment       string `json:"comment"`
	CreatedAt     string `json:"createdAt"`
}

// Trainer - модель тренера
type Trainer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Specialty string `json:"specialty"`
	Icon      string `json:"icon"`
}

// LoginRequest - запрос на вход
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest - запрос на регистрацию
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Response - стандартный ответ API
type Response struct {
	Token string      `json:"token,omitempty"`
	User  interface{} `json:"user,omitempty"`
	Error string      `json:"error,omitempty"`
}

// AvailableSlotsResponse - ответ с количеством свободных мест
type AvailableSlotsResponse struct {
	Available int `json:"available"`
}

// UpdateStatusRequest - запрос на обновление статуса бронирования
type UpdateStatusRequest struct {
	Status string `json:"status"`
}
