FROM golang:1.24.2-alpine AS builder

RUN mkdir /build
WORKDIR /build

COPY ./go.* .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w" \
    -o /build/bin/processor \
    ./cmd/processor/main.go

# Финальный образ с Docker CLI
FROM docker:latest

# Устанавливаем зависимости для запуска Go приложения
RUN apk add --no-cache ca-certificates

# Устанавливаем Docker клиент (уже есть в базовом образе)
# Копируем наше приложение
RUN mkdir /app
WORKDIR /app
COPY --from=builder /build/bin/processor /app/processor
COPY --from=builder /build/configs ./configs

# Убедимся, что можем запускать Docker внутри контейнера
# Для этого нужно монтировать docker.sock при запуске

CMD ["./processor"]