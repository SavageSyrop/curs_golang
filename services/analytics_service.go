package services

import (
	"banking-service/repositories"
	"fmt"
	"time"
)

// AnalyticsService управляет аналитикой финансовых операций
type AnalyticsService struct {
	TransactionRepo *repositories.TransactionRepository
}

// GetMonthlyIncomeAndExpenses возвращает доходы и расходы пользователя за указанный месяц
func (s *AnalyticsService) GetMonthlyIncomeAndExpenses(userID int, year, month int) (income float64, expenses float64, err error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // Первый день следующего месяца

	transactions, err := s.TransactionRepo.GetTransactionsByUserID(userID, startDate, endDate)
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось получить транзакции: %v", err)
	}

	for _, tx := range transactions {
		if tx.ToAccountID == userID { // Доход
			income += tx.Amount
		} else if tx.FromAccountID == userID { // Расход
			expenses += tx.Amount
		}
	}

	return income, expenses, nil
}

// PredictBalance прогнозирует баланс пользователя через N дней
func (s *AnalyticsService) PredictBalance(userID int, daysAhead int) (float64, error) {
	// Получаем текущий баланс пользователя
	balance, err := s.getTransactionBalance(userID)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить текущий баланс: %v", err)
	}

	// Получаем средний ежедневный доход и расход за последние 30 дней
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()
	income, expenses, err := s.GetMonthlyIncomeAndExpenses(userID, startDate.Year(), int(startDate.Month()))
	if err != nil {
		return 0, fmt.Errorf("не удалось рассчитать доходы/расходы: %v", err)
	}

	daysInMonth := endDate.Sub(startDate).Hours() / 24
	avgDailyIncome := income / daysInMonth
	avgDailyExpenses := expenses / daysInMonth

	// Прогнозируем баланс через N дней
	predictedBalance := balance + (avgDailyIncome-avgDailyExpenses)*float64(daysAhead)
	return predictedBalance, nil
}

// getTransactionBalance вычисляет текущий баланс пользователя на основе транзакций
func (s *AnalyticsService) getTransactionBalance(userID int) (float64, error) {
	transactions, err := s.TransactionRepo.GetTransactionsByUserID(userID, time.Time{}, time.Now())
	if err != nil {
		return 0, fmt.Errorf("не удалось получить транзакции: %v", err)
	}

	var balance float64
	for _, tx := range transactions {
		if tx.ToAccountID == userID { // Доход
			balance += tx.Amount
		} else if tx.FromAccountID == userID { // Расход
			balance -= tx.Amount
		}
	}

	return balance, nil
}
