package services

import (
	"banking-service/models"
	"banking-service/repositories"
	"errors"
	"fmt"
)

// AcountService управляет операциями со счетами
type AcountService struct {
	AccountRepo *repositories.AccountRepository
}

// CreateAccount создает новый счет для пользователя
func (s *AcountService) CreateAccount(userID int) (*models.Account, error) {
	account, err := s.AccountRepo.CreateAccount(userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать счет: %v", err)
	}
	return account, nil
}

// GetUserAccounts получает все счета пользователя
func (s *AcountService) GetUserAccounts(userID int) ([]models.Account, error) {
	accounts, err := s.AccountRepo.GetAccountsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить счета пользователя: %v", err)
	}
	return accounts, nil
}

// Deposit пополняет баланс счета
func (s *AcountService) Deposit(accountID int, amount float64) error {
	if amount <= 0 {
		return errors.New("сумма пополнения должна быть больше нуля")
	}

	err := s.AccountRepo.UpdateBalance(accountID, amount)
	if err != nil {
		return fmt.Errorf("не удалось пополнить баланс: %v", err)
	}
	return nil
}

// Withdraw списывает средства со счета
func (s *AcountService) Withdraw(accountID int, amount float64) error {
	if amount <= 0 {
		return errors.New("сумма списания должна быть больше нуля")
	}

	// Проверяем баланс перед списанием
	account, err := s.AccountRepo.GetAccountByID(accountID)
	if err != nil {
		return fmt.Errorf("не удалось найти счет: %v", err)
	}
	if account.Balance < amount {
		return errors.New("недостаточно средств на счете")
	}

	err = s.AccountRepo.UpdateBalance(accountID, -amount)
	if err != nil {
		return fmt.Errorf("не удалось списать средства: %v", err)
	}
	return nil
}
