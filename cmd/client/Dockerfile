# Stage 1: Build
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка клиента
RUN go build -ldflags="-s -w" -o /app/bin/client ./cmd/client/main.go

# Проверка содержимого /app/bin/
RUN ls -l /app/bin/

# Stage 2: Run
FROM alpine:latest

# Создаем рабочую директорию
WORKDIR /app

# Устанавливаем необходимые пакеты (включая `libc6-compat` для совместимости)
RUN apk add --no-cache libc6-compat

# Копируем собранный бинарный файл из первого этапа
COPY --from=builder /app/bin/client /app/client
COPY --from=builder /app/config/config.json /app/config/config.json

# Устанавливаем права на выполнение
RUN chmod +x /app/client

# Определяем точку входа
ENTRYPOINT ["/app/client"]
