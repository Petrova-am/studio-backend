package booking

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourusername/studio-backend/internal/models"
)

// Мок для репозитория
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(booking *models.Booking) (int, error) {
	args := m.Called(booking)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetByUserID(userID int) ([]models.Booking, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Booking), args.Error(1)
}

func (m *MockRepository) UpdateStatus(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockRepository) GetAvailableSlots(date, time string, trainerID int) (int, error) {
	args := m.Called(date, time, trainerID)
	return args.Int(0), args.Error(1)
}

// ========== ТЕСТЫ ==========

func TestCreateBooking_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	booking := &models.Booking{
		UserID:    1,
		Service:   "yoga",
		TrainerID: 1,
		Date:      "2026-05-20",
		Time:      "10:00",
		Amount:    1500,
	}

	mockRepo.On("Create", booking).Return(1, nil)

	service := NewService(mockRepo)
	err := service.CreateBooking(booking)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateBooking_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	booking := &models.Booking{}

	mockRepo.On("Create", booking).Return(0, errors.New("db error"))

	service := NewService(mockRepo)
	err := service.CreateBooking(booking)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUserBookings_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	expectedBookings := []models.Booking{
		{ID: 1, Service: "yoga", Status: "pending"},
		{ID: 2, Service: "stretch", Status: "confirmed"},
	}

	mockRepo.On("GetByUserID", 1).Return(expectedBookings, nil)

	service := NewService(mockRepo)
	bookings, err := service.GetUserBookings(1)

	assert.NoError(t, err)
	assert.Len(t, bookings, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetUserBookings_Error(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("GetByUserID", 1).Return(nil, errors.New("db error"))

	service := NewService(mockRepo)
	bookings, err := service.GetUserBookings(1)

	assert.Error(t, err)
	assert.Nil(t, bookings)
	mockRepo.AssertExpectations(t)
}

func TestConfirmBooking_Success(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("UpdateStatus", 1, "confirmed").Return(nil)

	service := NewService(mockRepo)
	err := service.ConfirmBooking(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestConfirmBooking_Error(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("UpdateStatus", 1, "confirmed").Return(errors.New("db error"))

	service := NewService(mockRepo)
	err := service.ConfirmBooking(1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCancelBooking_Success(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("UpdateStatus", 1, "cancelled").Return(nil)

	service := NewService(mockRepo)
	err := service.CancelBooking(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAvailableSlots_Success(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("GetAvailableSlots", "2026-05-20", "10:00", 1).Return(5, nil)

	service := NewService(mockRepo)
	slots, err := service.GetAvailableSlots("2026-05-20", "10:00", 1)

	assert.NoError(t, err)
	assert.Equal(t, 5, slots)
	mockRepo.AssertExpectations(t)
}

func TestGetAvailableSlots_Error(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("GetAvailableSlots", "2026-05-20", "10:00", 1).Return(0, errors.New("db error"))

	service := NewService(mockRepo)
	slots, err := service.GetAvailableSlots("2026-05-20", "10:00", 1)

	assert.Error(t, err)
	assert.Equal(t, 0, slots)
	mockRepo.AssertExpectations(t)
}
