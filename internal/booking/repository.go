package booking

import (
	"database/sql"

	"github.com/yourusername/studio-backend/internal/models"
)

type Repository interface {
	Create(booking *models.Booking) (int, error)
	GetByUserID(userID int) ([]models.Booking, error)
	UpdateStatus(id int, status string) error
	GetAvailableSlots(date, time string, trainerID int) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(booking *models.Booking) (int, error) {
	query := `INSERT INTO bookings (user_id, service, trainer_id, date, time, status, amount) 
              VALUES ($1, $2, $3, $4, $5, 'pending', $6) RETURNING id`
	err := r.db.QueryRow(query, booking.UserID, booking.Service, booking.TrainerID,
		booking.Date, booking.Time, booking.Amount).Scan(&booking.ID)
	return booking.ID, err
}

func (r *repository) GetByUserID(userID int) ([]models.Booking, error) {
	rows, err := r.db.Query("SELECT id, service, trainer_id, date, time, status FROM bookings WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		rows.Scan(&b.ID, &b.Service, &b.TrainerID, &b.Date, &b.Time, &b.Status)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *repository) UpdateStatus(id int, status string) error {
	_, err := r.db.Exec("UPDATE bookings SET status = $1 WHERE id = $2", status, id)
	return err
}

func (r *repository) GetAvailableSlots(date, time string, trainerID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM bookings WHERE date = $1 AND time = $2 AND trainer_id = $3"
	err := r.db.QueryRow(query, date, time, trainerID).Scan(&count)
	return 8 - count, err
}
