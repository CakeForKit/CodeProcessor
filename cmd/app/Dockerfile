# Стадия сборки (builder) - это "временный слой"
FROM golang:1.24.2-alpine AS builder

RUN apk add git bash && \
    go install github.com/swaggo/swag/cmd/swag@latest

RUN mkdir /build
WORKDIR /build

COPY ./go.* .
RUN go mod download

COPY . .

# Генерация Swagger документации
RUN /go/bin/swag init -g ./cmd/app/main.go --output ./docs


# Сборка Go приложения
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w" \
    -o /build/bin/code-processor-app \
    ./cmd/app/main.go

# Финальная стадия - Новый базовый слой
FROM alpine:latest  

RUN mkdir /app
WORKDIR /app
COPY --from=builder /build/bin/code-processor-app .
COPY --from=builder /build/configs ./configs
# COPY --from=builder /build/migrations ./migrations


RUN apk add --no-cache ca-certificates tzdata

# Порт, который слушает приложение
EXPOSE 8000

CMD ["./code-processor-app"]


