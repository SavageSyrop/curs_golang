package services

import (
	"banking-service/repositories"
	"errors"
	"fmt"
)

// TransferService управляет операциями перевода средств между счетами
type TransferService struct {
	TransactionRepo *repositories.TransactionRepository
}

// Transfer выполняет перевод средств между счетами
func (s *TransferService) Transfer(fromAccountID, toAccountID int, amount float64) error {
	if fromAccountID == toAccountID {
		return errors.New("невозможно перевести средства на тот же счет")
	}

	if amount <= 0 {
		return errors.New("сумма перевода должна быть больше нуля")
	}

	// Выполняем перевод средств
	err := s.TransactionRepo.TransferFunds(fromAccountID, toAccountID, amount)
	if err != nil {
		return fmt.Errorf("не удалось выполнить перевод: %v", err)
	}

	return nil
}
