# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder stage
COPY --from=builder /app/main .
# Копируем конфиги
COPY --from=builder /app/configs ./configs
# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Экспортируем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]