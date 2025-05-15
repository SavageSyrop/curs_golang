package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"banking-service/services"
)

// AccountHandler управляет HTTP-запросами, связанными со счетами
type AccountHandler struct {
	AccountService *services.AcountService
	CbrService     *services.CBRService
	EmailService   *services.EmailService
}

// CreateAccount обрабатывает создание нового счета
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста (предполагается, что он добавлен middleware)
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	// Создаем счет через сервис
	account, err := h.AccountService.CreateAccount(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем созданный счет в ответе
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// GetUserAccounts обрабатывает получение всех счетов пользователя
func (h *AccountHandler) GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	// Получаем счета пользователя через сервис
	accounts, err := h.AccountService.GetUserAccounts(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем список счетов в ответе
	json.NewEncoder(w).Encode(accounts)
}

func (h *AccountHandler) GetCbrInfo(w http.ResponseWriter, r *http.Request) {
	rate, err := h.CbrService.GetCentralBankRate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rate)
}

func (h *AccountHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email  string  `json:"email"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	err := h.EmailService.SendPaymentNotification(input.Email, input.Amount)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
