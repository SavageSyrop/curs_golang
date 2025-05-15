package services

import (
	"banking-service/models"
	"banking-service/repositories"
	"errors"
	"fmt"
	"math"
	"time"
)

// CreditService управляет операциями с кредитами
type CreditService struct {
	CreditRepo          *repositories.CreditRepository
	PaymentScheduleRepo *repositories.PaymentScheduleRepository
}

// CreateCredit создает новый кредит для пользователя
func (s *CreditService) CreateCredit(userID int, amount, rate float64, term int) error {
	if amount <= 0 || rate <= 0 || term <= 0 {
		return errors.New("сумма, ставка и срок кредита должны быть больше нуля")
	}

	// Создаем запись о кредите в базе данных
	err := s.CreditRepo.CreateCredit(userID, amount, rate, term)
	if err != nil {
		return fmt.Errorf("не удалось создать кредит: %v", err)
	}

	// Генерируем график платежей
	payments, err := s.GeneratePaymentSchedule(amount, rate, term)
	if err != nil {
		return fmt.Errorf("не удалось сгенерировать график платежей: %v", err)
	}

	// Получаем ID последнего созданного кредита
	creditID, err := s.CreditRepo.GetLastCreditIDByUserID(userID)
	if err != nil {
		return fmt.Errorf("не удалось получить ID кредита: %v", err)
	}

	// Сохраняем график платежей в базе данных
	err = s.PaymentScheduleRepo.CreatePaymentSchedule(creditID, payments)
	if err != nil {
		return fmt.Errorf("не удалось сохранить график платежей: %v", err)
	}

	return nil
}

// GeneratePaymentSchedule генерирует график аннуитетных платежей
func (s *CreditService) GeneratePaymentSchedule(amount, rate float64, term int) ([]models.PaymentSchedule, error) {
	monthlyRate := rate / 12 / 100 // Месячная процентная ставка
	annuityFactor := (monthlyRate * math.Pow(1+monthlyRate, float64(term))) /
		(math.Pow(1+monthlyRate, float64(term)) - 1)
	monthlyPayment := amount * annuityFactor

	var payments []models.PaymentSchedule
	remainingAmount := amount

	for i := 1; i <= term; i++ {
		interestPayment := remainingAmount * monthlyRate
		principalPayment := monthlyPayment - interestPayment
		remainingAmount -= principalPayment

		paymentDate := time.Now().AddDate(0, i, 0) // Дата платежа через i месяцев
		payments = append(payments, models.PaymentSchedule{
			PaymentDate: paymentDate,
			Amount:      monthlyPayment,
			Status:      "pending",
		})
	}

	return payments, nil
}

// ProcessOverduePayments обрабатывает просроченные платежи
func (s *CreditService) ProcessOverduePayments() {
	overduePayments, err := s.PaymentScheduleRepo.GetOverduePayments()
	if err != nil {
		fmt.Printf("Ошибка при получении просроченных платежей: %v\n", err)
		return
	}

	for _, payment := range overduePayments {
		// Начисляем штраф (например, +10% к сумме)
		newAmount := payment.Amount * 1.10

		// Обновляем сумму и статус платежа
		err := s.PaymentScheduleRepo.UpdatePaymentStatus(payment.ID, "overdue")
		if err != nil {
			fmt.Printf("Ошибка при обновлении статуса платежа: %v\n", err)
			continue
		}

		fmt.Printf("Просроченный платеж ID=%d обработан. Новая сумма: %.2f\n", payment.ID, newAmount)
	}
}

// GetPaymentScheduleByCreditID получает график платежей для конкретного кредита
func (s *CreditService) GetPaymentScheduleByCreditID(creditID int) ([]models.PaymentSchedule, error) {
	return s.PaymentScheduleRepo.GetPaymentScheduleByCreditID(creditID)
}
