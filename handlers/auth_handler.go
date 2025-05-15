package handlers

import (
	"banking-service/services"
	"encoding/json"
	"net/http"
)

// AuthHandler обрабатывает запросы, связанные с аутентификацией
type AuthHandler struct {
	AuthService *services.AuthService
}

// Register обрабатывает регистрацию пользователя
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := h.AuthService.Register(input.Username, input.Email, input.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь успешно зарегистрирован"))
}

// Login обрабатывает логин пользователя
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if jwt, err := h.AuthService.Login(input.Email, input.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.Header().Set("Authorization", jwt)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Авторизационный хедер в ответе"))
	}
}
