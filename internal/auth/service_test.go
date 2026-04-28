package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByEmail(email string) (*User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(email, passwordHash, name string) (int, error) {
	args := m.Called(email, passwordHash, name)
	return args.Int(0), args.Error(1)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	// Генерируем правильный хеш для пароля "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	mockRepo.On("GetUserByEmail", "test@test.com").Return(&User{
		ID:           1,
		Email:        "test@test.com",
		PasswordHash: string(hashedPassword),
		Name:         "Test User",
		Role:         "user",
	}, nil)

	service := NewAuthService(mockRepo, []byte("secret"))

	user, token, err := service.Login("test@test.com", "password")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserByEmail", "notfound@test.com").Return(nil, errors.New("not found"))

	service := NewAuthService(mockRepo, []byte("secret"))

	user, token, err := service.Login("notfound@test.com", "password")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("CreateUser", "new@test.com", mock.AnythingOfType("string"), "New User").Return(1, nil)

	service := NewAuthService(mockRepo, []byte("secret"))

	user, token, err := service.Register("new@test.com", "password", "New User")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestRegister_UserExists(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("CreateUser", "exists@test.com", mock.AnythingOfType("string"), "Exists").Return(0, errors.New("user already exists"))

	service := NewAuthService(mockRepo, []byte("secret"))

	user, token, err := service.Register("exists@test.com", "password", "Exists")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}
