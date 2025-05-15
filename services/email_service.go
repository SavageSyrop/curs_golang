package services

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/go-mail/mail/v2"
)

// EmailService управляет отправкой email-уведомлений
type EmailService struct{}

// SendPaymentNotification отправляет уведомление о платеже пользователю
func (s *EmailService) SendPaymentNotification(email string, amount float64) error {
	// Создание нового email-сообщения
	m := mail.NewMessage()
	m.SetHeader("From", "noreply@example.com") // Отправитель
	m.SetHeader("To", email)                   // Получатель
	m.SetHeader("Subject", "Уведомление о платеже")
	m.SetBody("text/html", fmt.Sprintf(`
        <h1>Спасибо за оплату!</h1>
        <p>Сумма: <strong>%.2f RUB</strong></p>
        <small>Это автоматическое уведомление</small>
    `, amount))

	// Настройка SMTP-клиента
	d := mail.NewDialer("smtp.example.com", 587, "user@example.com", "password")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: false} // Проверка сертификата

	// Отправка письма
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Ошибка при отправке email: %v", err)
		return fmt.Errorf("не удалось отправить email")
	}

	log.Printf("Email успешно отправлен на адрес: %s", email)
	return nil
}
