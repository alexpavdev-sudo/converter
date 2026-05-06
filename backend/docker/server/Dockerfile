# syntax=docker/dockerfile:1
# ---- общий этап сборки зависимостей ----
FROM golang:1.26.1-alpine AS deps
ARG GID=1000
ARG UID=1000
RUN addgroup -g ${GID} appgroup && \
    adduser -D -u ${UID} -G appgroup appuser
WORKDIR /home/appuser/app
COPY go.mod go.sum ./
RUN go mod download && go mod tidy
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.27.0
RUN apk add --no-cache gcc musl-dev

# ---- этап сборки бинарника для продакшена ----
FROM deps AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# ---- продакшен образ (минимальный) ----
FROM alpine:latest AS prod
ARG GID=1000
ARG UID=1000
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache tzdata
RUN addgroup -g ${GID} appgroup && \
    adduser -D -u ${UID} -G appgroup appuser
WORKDIR /home/appuser/app
COPY --from=builder /home/appuser/app/main .
RUN chown -R appuser:appgroup /home/appuser
USER appuser
CMD ["./main"]

# ---- образ для разработки (с live reload) ----
FROM deps AS dev
RUN go install github.com/air-verse/air@latest # Устанавливаем air для автоматической перезагрузки
RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY . .
COPY docker/server/entrypoint.sh /entrypoint.sh
RUN chown -R appuser:appgroup /home/appuser
RUN chown -R appuser:appgroup /entrypoint.sh
RUN chmod +x /entrypoint.sh
RUN go mod tidy
USER appuser
CMD ["air", "-c", "docker/server/.air.toml"]

# ---- образ для тестирования ----
FROM deps AS test
RUN apk add --no-cache bash ffmpeg
COPY --chown=appuser:appgroup . .
ENTRYPOINT ["sh", "-c", "go mod download && go mod tidy && go mod verify && exec go test \"$@\""]
CMD ["-v", "./..."]