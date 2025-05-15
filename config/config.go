package config

import (
	"fmt"
	"os"
)

// Config хранит конфигурацию приложения
type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret    string
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	CBRSOAPURL   string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	return &Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "yourpassword"),
		DBName:       getEnv("DB_NAME", "banking_service"),
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:     getIntEnv("SMTP_PORT", 587),
		SMTPUser:     getEnv("SMTP_USER", "user@example.com"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "password"),
		CBRSOAPURL:   getEnv("CBR_SOAP_URL", "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx"),
	}
}

// Вспомогательная функция для получения строковых переменных окружения
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Вспомогательная функция для получения целочисленных переменных окружения
func getIntEnv(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err == nil {
			return intValue
		}
	}
	return fallback
}
