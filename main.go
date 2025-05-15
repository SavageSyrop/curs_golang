package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"banking-service/config"
	"banking-service/handlers"
	"banking-service/middleware"
	"banking-service/repositories"
	"banking-service/services"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация подключения к базе данных
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	defer db.Close()

	// Инициализация репозиториев
	transactionRepo := &repositories.TransactionRepository{DB: db}
	userRepo := &repositories.UserRepository{DB: db}
	acountRepo := &repositories.AccountRepository{DB: db}
	cardRepo := &repositories.CardRepository{DB: db}
	creditRepo := &repositories.CreditRepository{DB: db}
	paymentRepo := &repositories.PaymentScheduleRepository{DB: db}

	// Инициализация сервисов
	transferService := &services.TransferService{TransactionRepo: transactionRepo}
	authService := &services.AuthService{UserRepo: userRepo}
	accountService := &services.AcountService{AccountRepo: acountRepo}
	analyticsService := &services.AnalyticsService{TransactionRepo: transactionRepo}
	cardService := &services.CardService{CardRepo: cardRepo}
	creditService := &services.CreditService{CreditRepo: creditRepo, PaymentScheduleRepo: paymentRepo}
	cbrService := &services.CBRService{}
	smtpService := &services.EmailService{}

	// Инициализация обработчиков
	transferHandler := &handlers.TransferHandler{TransferService: transferService}
	authHandler := &handlers.AuthHandler{AuthService: authService}
	analyticsHandler := &handlers.AnalyticsHandler{AnalyticsService: analyticsService}
	accountHandler := &handlers.AccountHandler{AccountService: accountService, CbrService: cbrService, EmailService: smtpService}
	cardHandler := &handlers.CardHandler{CardService: cardService}
	creditHandler := &handlers.CreditHandler{CreditService: creditService}

	// Настройка маршрутизатора
	r := mux.NewRouter()

	noAuthRouter := r.PathPrefix("/").Subrouter()
	noAuthRouter.Use()

	// Защищенные эндпоинты (требуют JWT-аутентификацию)
	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.AuthMiddleware)

	noAuthRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
	noAuthRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	noAuthRouter.HandleFunc("/cbr", accountHandler.GetCbrInfo).Methods("GET")

	authRouter.HandleFunc("/transfer", transferHandler.Transfer).Methods("POST")

	authRouter.HandleFunc("/wallet", accountHandler.CreateAccount).Methods("POST")
	authRouter.HandleFunc("/wallet", accountHandler.GetUserAccounts).Methods("GET")
	authRouter.HandleFunc("/email", accountHandler.SendEmail).Methods("POST")

	authRouter.HandleFunc("/vcard", cardHandler.GenerateCard).Methods("POST")
	authRouter.HandleFunc("/vcard", cardHandler.GetCards).Methods("GET")

	authRouter.HandleFunc("/credit", creditHandler.CreateCredit).Methods("POST")
	authRouter.HandleFunc("/credit", creditHandler.GetPaymentSchedule).Methods("GET")

	authRouter.HandleFunc("/balance", analyticsHandler.PredictBalance).Methods("POST")
	authRouter.HandleFunc("/balance", analyticsHandler.GetMonthlyAnalytics).Methods("GET")

	// Запуск сервера
	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
