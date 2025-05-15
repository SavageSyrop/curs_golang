package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"banking-service/services"
)

// TransferHandler управляет HTTP-запросами, связанными с переводами средств
type TransferHandler struct {
	TransferService *services.TransferService
}

// Transfer обрабатывает запрос на перевод средств между счетами
func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста (добавлен middleware)
	userIDStr := r.Context().Value("userID").(string)
	fromAccountID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	// Декодируем данные из тела запроса
	var input struct {
		ToAccountID int     `json:"to_account_id"`
		Amount      float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверяем входные данные
	if input.ToAccountID == fromAccountID {
		http.Error(w, "Нельзя перевести средства на тот же счет", http.StatusBadRequest)
		return
	}
	if input.Amount <= 0 {
		http.Error(w, "Сумма перевода должна быть больше нуля", http.StatusBadRequest)
		return
	}

	// Выполняем перевод через сервис
	err = h.TransferService.Transfer(fromAccountID, input.ToAccountID, input.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Перевод успешно выполнен"))
}
