FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Финальный образ
FROM alpine:3.18

WORKDIR /app

# Копируем бинарный файл из образа-сборщика
COPY --from=builder /app/main .

# Указываем, что порт 8080 будет открыт для контейнера
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]