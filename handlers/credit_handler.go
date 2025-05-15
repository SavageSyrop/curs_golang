package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"banking-service/services"
)

// CreditHandler управляет HTTP-запросами, связанными с кредитами
type CreditHandler struct {
	CreditService *services.CreditService
}

// CreateCredit обрабатывает запрос на создание нового кредита
func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста (добавлен middleware)
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	// Декодируем данные из тела запроса
	var input struct {
		Amount float64 `json:"amount"`
		Rate   float64 `json:"rate"`
		Term   int     `json:"term"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверяем входные данные
	if input.Amount <= 0 || input.Rate <= 0 || input.Term <= 0 {
		http.Error(w, "Сумма, ставка и срок кредита должны быть больше нуля", http.StatusBadRequest)
		return
	}

	// Создаем кредит через сервис
	err = h.CreditService.CreateCredit(userID, input.Amount, input.Rate, input.Term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Кредит успешно создан"))
}

// GetPaymentSchedule обрабатывает запрос на получение графика платежей по кредиту
func (h *CreditHandler) GetPaymentSchedule(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userIDStr := r.Context().Value("userID").(string)
	_, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	// Извлекаем параметр creditID из строки запроса
	creditID, err := strconv.Atoi(r.URL.Query().Get("credit_id"))
	if err != nil || creditID <= 0 {
		http.Error(w, "Неверный ID кредита", http.StatusBadRequest)
		return
	}

	// Проверяем, что кредит принадлежит пользователю (опционально)
	// Например, вызываем метод CreditRepo.GetCreditByID(creditID) и проверяем поле UserID

	// Получаем график платежей через сервис
	schedule, err := h.CreditService.GetPaymentScheduleByCreditID(creditID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем график платежей в ответе
	json.NewEncoder(w).Encode(schedule)
}
