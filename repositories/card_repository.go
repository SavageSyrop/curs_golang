package repositories

import (
	"banking-service/models"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// CardRepository управляет операциями с картами в базе данных
type CardRepository struct {
	DB *sql.DB
}

// CreateCard создает новую карту для пользователя
func (r *CardRepository) CreateCard(userID int, number, expiry, cvv string) (*models.Card, error) {
	// Шифрование номера и срока действия через PGP (заглушка)
	numberEncrypted := "encrypted_" + number
	expiryEncrypted := "encrypted_" + expiry

	// Хеширование CVV через bcrypt
	cvvHashed, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("не удалось захешировать CVV: %v", err)
	}

	// Вычисление HMAC для номера карты
	hmac := models.ComputeHMAC(number, []byte("secret_key"))

	query := `INSERT INTO cards (user_id, number_encrypted, expiry_encrypted, cvv_hashed, hmac) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var card models.Card
	err = r.DB.QueryRow(query, userID, numberEncrypted, expiryEncrypted, string(cvvHashed), hmac).
		Scan(&card.ID)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать карту: %v", err)
	}
	card.UserID = userID
	card.NumberEncrypted = numberEncrypted
	card.ExpiryEncrypted = expiryEncrypted
	card.CVVHashed = string(cvvHashed)
	card.HMAC = hmac
	return &card, nil
}
