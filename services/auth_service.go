package services

import (
	"banking-service/config"
	"banking-service/models"
	"banking-service/repositories"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// AuthService управляет процессами регистрации и аутентификации
type AuthService struct {
	UserRepo *repositories.UserRepository
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(username, email, password string) error {
	user := &models.User{
		Username: username,
		Email:    email,
	}
	if err := user.HashPassword(password); err != nil {
		return errors.New("не удалось захешировать пароль")
	}

	existingUser, _ := s.UserRepo.GetUserByEmail(email)
	if existingUser != nil {
		return errors.New("email уже используется")
	}

	return s.UserRepo.CreateUser(user)
}

// Login аутентифицирует пользователя и возвращает JWT-токен
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("неверные данные для входа")
	}

	if !user.CheckPassword(password) {
		return "", errors.New("неверные данные для входа")
	}

	// Секретный ключ для подписи токена
	var mySigningKey = []byte(config.LoadConfig().JWTSecret)

	// Создание нового токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Время истечения токена
	})

	// Подпись токена
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("Ошибка при подписывании токена:", err)
		return "", nil
	}

	return "Bearer " + tokenString, nil
}
