package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Card представляет банковскую карту
type Card struct {
	ID              int    `json:"id"`
	UserID          int    `json:"-"`
	NumberEncrypted string `json:"-"` // Зашифрован PGP
	ExpiryEncrypted string `json:"-"` // Зашифрован PGP
	CVVHashed       string `json:"-"` // Хеширован bcrypt
	HMAC            string `json:"hmac"`
}

// ComputeHMAC вычисляет HMAC для номера карты
func ComputeHMAC(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
