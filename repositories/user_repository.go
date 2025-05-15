package repositories

import (
	"banking-service/models"
	"database/sql"
	"fmt"
)

// UserRepository управляет операциями с пользователями в базе данных
type UserRepository struct {
	DB *sql.DB
}

// CreateUser создает нового пользователя в базе данных
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("не удалось создать пользователя: %v", err)
	}
	return nil
}

// GetUserByEmail получает пользователя по email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash FROM users WHERE email = $1`
	row := r.DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось получить пользователя: %v", err)
	}
	return &user, nil
}
