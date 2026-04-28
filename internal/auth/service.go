package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService - интерфейс сервиса авторизации
type AuthService interface {
	Login(email, password string) (*User, string, error)
	Register(email, password, name string) (*User, string, error)
}

type authService struct {
	repo      UserRepository
	jwtSecret []byte
}

func NewAuthService(repo UserRepository, jwtSecret []byte) AuthService {
	return &authService{repo: repo, jwtSecret: jwtSecret}
}

func (s *authService) Login(email, password string) (*User, string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, _ := s.generateJWT(user.ID)
	return user, token, nil
}

func (s *authService) Register(email, password, name string) (*User, string, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	id, err := s.repo.CreateUser(email, string(hashedPassword), name)
	if err != nil {
		return nil, "", errors.New("user already exists")
	}

	token, _ := s.generateJWT(id)
	return &User{ID: id, Email: email, Name: name, Role: "user"}, token, nil
}

func (s *authService) generateJWT(userID int) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   string(rune(userID)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
