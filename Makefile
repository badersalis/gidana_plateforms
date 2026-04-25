run:
	go run ./cmd/main.go

build:
	go build -o gidana_api ./cmd/main.go

tidy:
	go mod tidy

migrate:
	go run ./cmd/main.go

.PHONY: run build tidy migrate
