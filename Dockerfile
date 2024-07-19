# Указываем базовый образ
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Сборка приложения
RUN go build -o main .

# Указываем базовый образ для выполнения
FROM debian:bullseye-slim

# Копируем собранное приложение из предыдущего шага
COPY --from=builder /app/main /app/main

# Открываем порт, который использует ваше приложение
EXPOSE 4000

# Команда для запуска приложения
CMD ["/app/main"]
