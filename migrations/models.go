package models

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

type Booking struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	Service   string `json:"service"`
	TrainerID int    `json:"trainerId"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Status    string `json:"status"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"createdAt"`
}

type Trainer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Specialty string `json:"specialty"`
	Icon      string `json:"icon"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Response struct {
	Token string      `json:"token,omitempty"`
	User  interface{} `json:"user,omitempty"`
	Error string      `json:"error,omitempty"`
}
