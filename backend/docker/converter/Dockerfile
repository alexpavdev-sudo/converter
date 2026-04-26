# syntax=docker/dockerfile:1
# ---- общий этап сборки зависимостей ----
FROM golang:1.26.1-alpine AS deps
ARG GID=1000
ARG UID=1000
RUN addgroup -g ${GID} appgroup && \
    adduser -D -u ${UID} -G appgroup appuser
WORKDIR /home/appuser/app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.27.0

# ---- этап сборки бинарника для продакшена ----
FROM deps AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o converter ./cmd/converter

# ---- продакшен образ (минимальный) ----
FROM alpine:latest AS prod
ARG GID=1000
ARG UID=1000
RUN apk --no-cache add ca-certificates
RUN addgroup -g ${GID} appgroup && \
    adduser -D -u ${UID} -G appgroup appuser
WORKDIR /home/appuser/app
COPY --from=builder /home/appuser/app/converter .
RUN chown -R appuser:appgroup /home/appuser
USER appuser
CMD ["./converter"]

# ---- образ для разработки (с live reload) ----
FROM deps AS dev
RUN apk add --no-cache gcc musl-dev # компилятор C для Alpine
RUN go install github.com/air-verse/air@latest # Устанавливаем air для автоматической перезагрузки
RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY . .
COPY docker/converter/entrypoint.sh /entrypoint.sh
RUN chown -R appuser:appgroup /home/appuser
RUN chown -R appuser:appgroup /entrypoint.sh
RUN chmod +x /entrypoint.sh
USER appuser
CMD ["air", "-c", "docker/converter/.air.toml"]