package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/beevik/etree"
)

// CBRService управляет интеграцией с Центральным банком РФ через SOAP
type CBRService struct{}

// GetCentralBankRate получает ключевую ставку ЦБ РФ за последние 30 дней
func (s *CBRService) GetCentralBankRate() (float64, error) {
	// Формирование дат для запроса
	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	// Формирование SOAP-запроса
	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
        <soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
            <soap12:Body>
                <KeyRate xmlns="http://web.cbr.ru/">
                    <fromDate>%s</fromDate>
                    <ToDate>%s</ToDate>
                </KeyRate>
            </soap12:Body>
        </soap12:Envelope>`, fromDate, toDate)

	// Отправка SOAP-запроса
	resp, err := http.Post("https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx", "application/soap+xml; charset=utf-8", bytes.NewBuffer([]byte(soapRequest)))
	if err != nil {
		return 0, fmt.Errorf("ошибка при отправке SOAP-запроса: %v", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	// Парсинг XML-ответа
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(rawBody); err != nil {
		return 0, fmt.Errorf("ошибка при парсинге XML: %v", err)
	}

	// Поиск элемента с ключевой ставкой
	rateElement := doc.FindElement("//diffgram/KeyRate/KR/Rate")
	if rateElement == nil {
		return 0, fmt.Errorf("ключевая ставка не найдена в ответе")
	}

	// Преобразование значения в число
	var rate float64
	_, err = fmt.Sscanf(rateElement.Text(), "%f", &rate)
	if err != nil {
		return 0, fmt.Errorf("ошибка при преобразовании ставки: %v", err)
	}

	log.Printf("Получена ключевая ставка ЦБ РФ: %.2f%%", rate)
	return rate, nil
}
