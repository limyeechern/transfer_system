.PHONY: run build test fmt clean docker-build docker-up docker-down docker-logs

APP_NAME := transfer-system

run:
	go run ./...

build:
	go build -o bin/$(APP_NAME) .

test:
	go test ./...

fmt:
	go fmt ./...

clean:
	rm -rf bin

docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f
