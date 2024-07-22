# Stage 1: Build
FROM golang:1.22-alpine AS builder

# Установка рабочего каталога
WORKDIR /app

# Копирование файлов go.mod и go.sum
COPY go.mod go.sum ./

# Установка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Построение приложения
RUN go build -ldflags="-s -w" -o /app/bin/server ./cmd/message-service/main.go

# Stage 2: Run
FROM alpine:latest

# Установка рабочего каталога
WORKDIR /app

# Копирование построенного приложения из предыдущего этапа
COPY --from=builder /app/bin/server /app/server

# Копирование конфигурации файла 
COPY --from=builder /app/config/config.json /app/config/config.json

# Установка необходимых пакетов и утилиты dumb-init для обработки сигналов
RUN apk add --no-cache dumb-init libgcc libstdc++ ca-certificates

# Делаем файл исполняемым
RUN chmod +x /app/server

# Устанавливаем переменную окружения для режима prefork
ENV FIBER_PREFORK=true

# Открываем порт для приложения
EXPOSE 4000

# Используем dumb-init для управления процессами и сигналами
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# Команда по умолчанию
CMD ["/app/server"]
