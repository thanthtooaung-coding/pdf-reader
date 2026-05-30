FROM golang:1.22-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /pdf-reader-backend ./cmd

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /pdf-reader-backend .
RUN mkdir -p /app/storage

EXPOSE 8080
CMD ["./pdf-reader-backend"]
