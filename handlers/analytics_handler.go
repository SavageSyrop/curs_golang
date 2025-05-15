package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"banking-service/services"
)

// AnalyticsHandler управляет HTTP-запросами, связанными с аналитикой
type AnalyticsHandler struct {
	AnalyticsService *services.AnalyticsService
}

// GetMonthlyAnalytics обрабатывает запрос на получение доходов и расходов за месяц
func (h *AnalyticsHandler) GetMonthlyAnalytics(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста (добавлен middleware)
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	// Извлекаем параметры year и month из строки запроса
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil || year <= 0 {
		http.Error(w, "Неверный год", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil || month < 1 || month > 12 {
		http.Error(w, "Неверный месяц", http.StatusBadRequest)
		return
	}

	// Получаем доходы и расходы через сервис
	income, expenses, err := h.AnalyticsService.GetMonthlyIncomeAndExpenses(userID, year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем результат в формате JSON
	response := map[string]float64{
		"income":   income,
		"expenses": expenses,
	}
	json.NewEncoder(w).Encode(response)
}

// PredictBalance обрабатывает запрос на прогноз баланса
func (h *AnalyticsHandler) PredictBalance(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	// Извлекаем параметр daysAhead из строки запроса
	daysAhead, err := strconv.Atoi(r.URL.Query().Get("days"))
	if err != nil || daysAhead <= 0 {
		http.Error(w, "Количество дней должно быть больше нуля", http.StatusBadRequest)
		return
	}

	// Прогнозируем баланс через сервис
	predictedBalance, err := h.AnalyticsService.PredictBalance(userID, daysAhead)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем результат в формате JSON
	response := map[string]float64{
		"predicted_balance": predictedBalance,
	}
	json.NewEncoder(w).Encode(response)
}
