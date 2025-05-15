package repositories

import (
	"banking-service/models"
	"database/sql"
	"fmt"
)

// AccountRepository управляет операциями со счетами в базе данных
type AccountRepository struct {
	DB *sql.DB
}

// CreateAccount создает новый счет для пользователя
func (r *AccountRepository) CreateAccount(userID int) (*models.Account, error) {
	query := `INSERT INTO accounts (user_id, balance) VALUES ($1, $2) RETURNING id, balance`
	var account models.Account
	err := r.DB.QueryRow(query, userID, 0).Scan(&account.ID, &account.Balance)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать счет: %v", err)
	}
	account.UserID = userID
	return &account, nil
}

// GetAccountsByUserID получает все счета пользователя
func (r *AccountRepository) GetAccountsByUserID(userID int) ([]models.Account, error) {
	query := `SELECT id, user_id, balance FROM accounts WHERE user_id = $1`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить счета: %v", err)
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		if err := rows.Scan(&account.ID, &account.UserID, &account.Balance); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании счета: %v", err)
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// UpdateBalance обновляет баланс счета
func (r *AccountRepository) UpdateBalance(accountID int, amount float64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err := r.DB.Exec(query, amount, accountID)
	if err != nil {
		return fmt.Errorf("не удалось обновить баланс: %v", err)
	}
	return nil
}

// GetAccountByID получает информацию о счете по его ID
func (r *AccountRepository) GetAccountByID(accountID int) (*models.Account, error) {
	query := `SELECT id, user_id, balance FROM accounts WHERE id = $1`
	row := r.DB.QueryRow(query, accountID)

	var account models.Account
	err := row.Scan(&account.ID, &account.UserID, &account.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось получить счет: %v", err)
	}
	return &account, nil
}
