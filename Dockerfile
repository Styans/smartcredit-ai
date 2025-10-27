# Dockerfile

# --- Сборочный этап ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# === ИЗМЕНЕНИЕ: Указываем v1.62.0 ===
RUN go install github.com/air-verse/air@v1.62.0

# Копируем исходный код
COPY . .

# Собираем наше основное приложение
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/api

# --- Финальный этап ---
FROM alpine:3.19

WORKDIR /app

# Копируем наш скомпилированный бинарник
COPY --from=builder /app/server .

# Копируем 'air' из builder-образа
COPY --from=builder /go/bin/air /usr/local/bin/air

# Копируем .env
COPY .env .

EXPOSE 9090
CMD ["./server"]