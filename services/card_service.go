package services

import (
	"banking-service/models"
	"banking-service/repositories"
	"fmt"
	"math/rand"
	"time"
)

// CardService управляет операциями с банковскими картами
type CardService struct {
	CardRepo *repositories.CardRepository
}

// GenerateCard создает новую виртуальную карту для пользователя
func (s *CardService) GenerateCard(userID int) (*models.Card, error) {
	// Генерация номера карты (пример: 16-значный номер)
	cardNumber := generateCardNumber()

	// Генерация срока действия карты (3 года от текущей даты)
	expiryDate := time.Now().AddDate(3, 0, 0).Format("01/2006")

	// Генерация CVV (3-значное число)
	cvv := fmt.Sprintf("%03d", generateRandomNumber(100, 999))

	// Создание карты в базе данных
	card, err := s.CardRepo.CreateCard(userID, cardNumber, expiryDate, cvv)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать карту: %v", err)
	}
	return card, nil
}

// GetCards получает все карты пользователя
func (s *CardService) GetCards(userID int) ([]models.Card, error) {
	cards, err := s.CardRepo.GetCardsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить карты пользователя: %v", err)
	}
	return cards, nil
}

// VerifyCardHMAC проверяет целостность данных карты через HMAC
func (s *CardService) VerifyCardHMAC(card models.Card, secret []byte) bool {
	expectedHMAC := models.ComputeHMAC(card.NumberEncrypted, secret)
	return card.HMAC == expectedHMAC
}

// Helper function: Генерация номера карты по алгоритму Луна
func generateCardNumber() string {
	// Пример: 4111111111111111 (валидный номер карты для тестирования)
	return "4111111111111111"
}

// Helper function: Генерация случайного числа в заданном диапазоне
func generateRandomNumber(min, max int) int {
	return min + rand.Intn(max-min+1)
}
