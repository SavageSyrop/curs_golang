package repositories

import (
	"banking-service/models"
	"database/sql"
	"errors"
	"fmt"
)

// CreditRepository управляет операциями с кредитами в базе данных
type CreditRepository struct {
	DB *sql.DB
}

// CreateCredit создает новый кредит для пользователя
func (r *CreditRepository) CreateCredit(userID int, amount, rate float64, term int) error {
	query := `INSERT INTO credits (user_id, amount, rate, term) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, userID, amount, rate, term)
	if err != nil {
		return fmt.Errorf("не удалось создать кредит: %v", err)
	}
	return nil
}

// GetCreditsByUserID получает все кредиты пользователя
func (r *CreditRepository) GetCreditsByUserID(userID int) ([]models.Credit, error) {
	query := `SELECT id, user_id, amount, rate, term FROM credits WHERE user_id = $1`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить кредиты: %v", err)
	}
	defer rows.Close()

	var credits []models.Credit
	for rows.Next() {
		var credit models.Credit
		if err := rows.Scan(&credit.ID, &credit.UserID, &credit.Amount, &credit.Rate, &credit.Term); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании кредита: %v", err)
		}
		credits = append(credits, credit)
	}
	return credits, nil
}

// CreatePaymentSchedule создает график платежей для кредита
func (r *CreditRepository) CreatePaymentSchedule(creditID int, payments []models.PaymentSchedule) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %v", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO payment_schedules (credit_id, payment_date, amount, status) VALUES ($1, $2, $3, $4)`
	for _, payment := range payments {
		_, err := tx.Exec(query, creditID, payment.PaymentDate, payment.Amount, payment.Status)
		if err != nil {
			return fmt.Errorf("не удалось создать запись в графике платежей: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("не удалось завершить транзакцию: %v", err)
	}
	return nil
}

// GetPaymentScheduleByCreditID получает график платежей для конкретного кредита
func (r *CreditRepository) GetPaymentScheduleByCreditID(creditID int) ([]models.PaymentSchedule, error) {
	query := `SELECT id, credit_id, payment_date, amount, status FROM payment_schedules WHERE credit_id = $1`
	rows, err := r.DB.Query(query, creditID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить график платежей: %v", err)
	}
	defer rows.Close()

	var schedule []models.PaymentSchedule
	for rows.Next() {
		var payment models.PaymentSchedule
		if err := rows.Scan(&payment.ID, &payment.CreditID, &payment.PaymentDate, &payment.Amount, &payment.Status); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании графика платежей: %v", err)
		}
		schedule = append(schedule, payment)
	}
	return schedule, nil
}

// GetCardsByUserID получает все карты пользователя
func (r *CardRepository) GetCardsByUserID(userID int) ([]models.Card, error) {
	query := `SELECT id, user_id, number_encrypted, expiry_encrypted, cvv_hashed, hmac FROM cards WHERE user_id = $1`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить карты пользователя: %v", err)
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var card models.Card
		if err := rows.Scan(&card.ID, &card.UserID, &card.NumberEncrypted, &card.ExpiryEncrypted, &card.CVVHashed, &card.HMAC); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании карты: %v", err)
		}
		cards = append(cards, card)
	}
	return cards, nil
}

// GetLastCreditIDByUserID получает ID последнего созданного кредита пользователя
func (r *CreditRepository) GetLastCreditIDByUserID(userID int) (int, error) {
	query := `SELECT id FROM credits WHERE user_id = $1 ORDER BY id DESC LIMIT 1`
	row := r.DB.QueryRow(query, userID)

	var creditID int
	err := row.Scan(&creditID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("кредиты не найдены")
		}
		return 0, fmt.Errorf("не удалось получить ID кредита: %v", err)
	}
	return creditID, nil
}
