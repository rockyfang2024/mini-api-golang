# ---- Build stage ----
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Install gcc for cgo (required by mattn/go-sqlite3)
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o mini-api ./cmd/main.go

# ---- Runtime stage ----
FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/mini-api .
COPY config/app.yaml config/app.yaml

EXPOSE 8080
CMD ["./mini-api"]
