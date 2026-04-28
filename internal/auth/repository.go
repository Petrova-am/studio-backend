package auth

import "database/sql"

// User - модель пользователя
type User struct {
	ID           int
	Email        string
	PasswordHash string
	Name         string
	Role         string
}

// UserRepository - интерфейс для работы с пользователями
type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(email, passwordHash, name string) (int, error)
}

// userRepository - реальная реализация с БД
type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	query := "SELECT id, email, password_hash, name, role FROM users WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(email, passwordHash, name string) (int, error) {
	var id int
	query := "INSERT INTO users (email, password_hash, name, role) VALUES ($1, $2, $3, 'user') RETURNING id"
	err := r.db.QueryRow(query, email, passwordHash, name).Scan(&id)
	return id, err
}
