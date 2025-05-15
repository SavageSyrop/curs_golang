# Базовый образ с Go 1.24.2
FROM golang:1.24.2-alpine AS builder

# Установка зависимостей
RUN apk add --no-cache git

# Рабочая директория
WORKDIR /app

# Копирование модулей и кода
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o banking-service .

# Финальный образ
FROM alpine:latest

# Установка tzdata для работы с временными зонами
RUN apk add --no-cache tzdata

# Копирование бинарного файла
WORKDIR /root/
COPY --from=builder /app/banking-service .

# Команда запуска
CMD ["./banking-service"]