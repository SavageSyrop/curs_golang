package repositories

import (
	"banking-service/models"
	"database/sql"
	"fmt"
	"time"
)

// TransactionRepository управляет операциями с транзакциями в базе данных
type TransactionRepository struct {
	DB *sql.DB
}

// CreateTransaction создает новую транзакцию
func (r *TransactionRepository) CreateTransaction(fromAccountID, toAccountID int, amount float64) error {
	query := `
        INSERT INTO transactions (from_account_id, to_account_id, amount, created_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.DB.Exec(query, fromAccountID, toAccountID, amount, time.Now())
	if err != nil {
		return fmt.Errorf("не удалось создать транзакцию: %v", err)
	}
	return nil
}

// GetTransactionsByAccountID получает историю транзакций для конкретного счета
func (r *TransactionRepository) GetTransactionsByAccountID(accountID int) ([]models.Transaction, error) {
	query := `
        SELECT id, from_account_id, to_account_id, amount, created_at 
        FROM transactions 
        WHERE from_account_id = $1 OR to_account_id = $1
        ORDER BY created_at DESC
    `
	rows, err := r.DB.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить транзакции: %v", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.FromAccountID, &transaction.ToAccountID, &transaction.Amount, &transaction.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании транзакции: %v", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// TransferFunds выполняет перевод средств между счетами с использованием транзакций
func (r *TransactionRepository) TransferFunds(fromAccountID, toAccountID int, amount float64) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %v", err)
	}
	defer tx.Rollback()

	// Проверяем баланс отправителя
	var senderBalance float64
	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, fromAccountID).Scan(&senderBalance)
	if err != nil {
		return fmt.Errorf("не удалось получить баланс отправителя: %v", err)
	}
	if senderBalance < amount {
		return fmt.Errorf("недостаточно средств на счете отправителя")
	}

	// Списываем средства со счета отправителя
	_, err = tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, fromAccountID)
	if err != nil {
		return fmt.Errorf("не удалось списать средства со счета отправителя: %v", err)
	}

	// Зачисляем средства на счет получателя
	_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toAccountID)
	if err != nil {
		return fmt.Errorf("не удалось зачислить средства на счет получателя: %v", err)
	}

	// Создаем запись о транзакции
	_, err = tx.Exec(`
        INSERT INTO transactions (from_account_id, to_account_id, amount, created_at)
        VALUES ($1, $2, $3, $4)
    `, fromAccountID, toAccountID, amount, time.Now())
	if err != nil {
		return fmt.Errorf("не удалось создать запись о транзакции: %v", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("не удалось завершить транзакцию: %v", err)
	}
	return nil
}

// GetTransactionsByUserID получает транзакции пользователя за указанный период
func (r *TransactionRepository) GetTransactionsByUserID(userID int, startDate, endDate time.Time) ([]models.Transaction, error) {
	query := `
        SELECT id, from_account_id, to_account_id, amount, created_at 
        FROM transactions 
        WHERE (from_account_id = $1 OR to_account_id = $1) AND created_at >= $2 AND created_at < $3
        ORDER BY created_at DESC
    `
	rows, err := r.DB.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить транзакции: %v", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.FromAccountID, &transaction.ToAccountID, &transaction.Amount, &transaction.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании транзакции: %v", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
