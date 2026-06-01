FROM golang:1.24-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .
RUN swag init -g main.go -d cmd,internal/handler,internal/response -o docs --parseInternal --parseDependencyLevel 1
RUN CGO_ENABLED=0 GOOS=linux go build -o /pdf-reader-backend ./cmd

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /pdf-reader-backend .
RUN mkdir -p /app/storage

EXPOSE 8080
CMD ["./pdf-reader-backend"]
