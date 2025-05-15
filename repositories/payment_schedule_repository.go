package repositories

import (
	"banking-service/models"
	"database/sql"
	"fmt"
)

// PaymentScheduleRepository управляет операциями с графиком платежей в базе данных
type PaymentScheduleRepository struct {
	DB *sql.DB
}

// CreatePaymentSchedule создает график платежей для кредита
func (r *PaymentScheduleRepository) CreatePaymentSchedule(creditID int, payments []models.PaymentSchedule) error {
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

// UpdatePaymentStatus обновляет статус платежа (например, на "paid" или "overdue")
func (r *PaymentScheduleRepository) UpdatePaymentStatus(paymentID int, status string) error {
	query := `UPDATE payment_schedules SET status = $1 WHERE id = $2`
	_, err := r.DB.Exec(query, status, paymentID)
	if err != nil {
		return fmt.Errorf("не удалось обновить статус платежа: %v", err)
	}
	return nil
}

// GetOverduePayments получает все просроченные платежи
func (r *PaymentScheduleRepository) GetOverduePayments() ([]models.PaymentSchedule, error) {
	query := `
        SELECT id, credit_id, payment_date, amount, status 
        FROM payment_schedules 
        WHERE status = 'pending' AND payment_date < NOW()
    `
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить просроченные платежи: %v", err)
	}
	defer rows.Close()

	var overduePayments []models.PaymentSchedule
	for rows.Next() {
		var payment models.PaymentSchedule
		if err := rows.Scan(&payment.ID, &payment.CreditID, &payment.PaymentDate, &payment.Amount, &payment.Status); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании просроченного платежа: %v", err)
		}
		overduePayments = append(overduePayments, payment)
	}
	return overduePayments, nil
}

// GetPaymentScheduleByCreditID получает график платежей для конкретного кредита
func (r *PaymentScheduleRepository) GetPaymentScheduleByCreditID(creditID int) ([]models.PaymentSchedule, error) {
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
