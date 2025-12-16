# ---------- BUILD STAGE ----------
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# сборка worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o worker ./cmd/worker


# ---------- RUNTIME STAGE ----------
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/worker /app/worker

ENTRYPOINT ["/app/worker"]
