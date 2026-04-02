# syntax=docker/dockerfile:1

# ---- общий этап сборки зависимостей ----
FROM golang:1.26.1-alpine AS deps

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# ---- этап сборки бинарника для продакшена ----
FROM deps AS builder

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# ---- продакшен образ (минимальный) ----
FROM alpine:latest AS prod

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]

# ---- образ для разработки (с live reload) ----
FROM deps AS dev

# Устанавливаем air для автоматической перезагрузки
RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app
COPY . .
COPY .air.toml .
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

# Используем air с конфигурацией по умолчанию
CMD ["air", "-c", ".air.toml"]