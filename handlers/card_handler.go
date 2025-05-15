package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"banking-service/services"
)

// CardHandler управляет HTTP-запросами, связанными с банковскими картами
type CardHandler struct {
	CardService *services.CardService
}

// GenerateCard обрабатывает запрос на создание новой карты
func (h *CardHandler) GenerateCard(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста (добавлен middleware)
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	// Создаем карту через сервис
	card, err := h.CardService.GenerateCard(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем созданную карту в ответе
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

// GetCards обрабатывает запрос на получение всех карт пользователя
func (h *CardHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из контекста
	userIDStr := r.Context().Value("userID").(string)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	// Получаем карты пользователя через сервис
	cards, err := h.CardService.GetCards(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем список карт в ответе
	json.NewEncoder(w).Encode(cards)
}
