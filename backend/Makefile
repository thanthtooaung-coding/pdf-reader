.PHONY: run build test swagger

run:
	go run ./cmd

build:
	go build -o pdf-reader-backend ./cmd

test:
	go test ./... -count=1

swagger:
	swag init -g main.go -d cmd,internal/handler,internal/response -o docs --parseInternal --parseDependencyLevel 1
