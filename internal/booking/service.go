package booking

import (
	"github.com/yourusername/studio-backend/internal/models"
)

type Service interface {
	CreateBooking(booking *models.Booking) error
	GetUserBookings(userID int) ([]models.Booking, error)
	ConfirmBooking(id int) error
	CancelBooking(id int) error
	GetAvailableSlots(date, time string, trainerID int) (int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateBooking(booking *models.Booking) error {
	_, err := s.repo.Create(booking)
	return err
}

func (s *service) GetUserBookings(userID int) ([]models.Booking, error) {
	return s.repo.GetByUserID(userID)
}

func (s *service) ConfirmBooking(id int) error {
	return s.repo.UpdateStatus(id, "confirmed")
}

func (s *service) CancelBooking(id int) error {
	return s.repo.UpdateStatus(id, "cancelled")
}

func (s *service) GetAvailableSlots(date, time string, trainerID int) (int, error) {
	return s.repo.GetAvailableSlots(date, time, trainerID)
}
