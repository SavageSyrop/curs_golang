package models

// Credit представляет кредит пользователя
type Credit struct {
	ID     int     `json:"id"`
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
	Rate   float64 `json:"rate"`
	Term   int     `json:"term"`
}
