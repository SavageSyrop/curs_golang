package middleware

import (
	"banking-service/config"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware проверяет JWT-токен в заголовке Authorization
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		// Извлекаем токен из заголовка (убираем префикс "Bearer ")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Если префикс не был найден
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}

		var mySigningKey = []byte(config.LoadConfig().JWTSecret)

		// Парсинг токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверка метода подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неверный метод подписи")
			}
			return mySigningKey, nil
		})

		if err != nil {
			fmt.Println("Ошибка при парсинге токена:", err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println("Пользователь:", claims["username"])
			fmt.Println("Время истечения:", claims["exp"])
			str := strconv.FormatFloat(claims["sub"].(float64), 'f', -1, 64)
			ctx := context.WithValue(r.Context(), "userID", str)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			fmt.Println("Токен недействителен")
		}
	})
}
