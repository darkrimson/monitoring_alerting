# ---------- BUILD STAGE ----------
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# зависимости
COPY go.mod go.sum ./
RUN go mod download

# весь проект
COPY . .

# сборка API
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o api ./cmd/api


# ---------- RUNTIME STAGE ----------
FROM alpine:latest

WORKDIR /app

# сертификаты для https
RUN apk add --no-cache ca-certificates

# бинарь
COPY --from=builder /app/api /app/api

EXPOSE 8080

ENTRYPOINT ["/app/api"]
