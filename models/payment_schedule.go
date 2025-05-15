package models

import "time"

// PaymentSchedule представляет график платежей по кредиту
type PaymentSchedule struct {
	ID          int       `json:"id"`
	CreditID    int       `json:"credit_id"`
	PaymentDate time.Time `json:"payment_date"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"` // pending, paid, overdue
}
