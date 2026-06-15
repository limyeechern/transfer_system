.PHONY: run build test fmt clean docker-build docker-up docker-down docker-logs

APP_NAME := transfer-system

test:
	go test ./...

fmt:
	go fmt ./...

clean:
	rm -rf bin

build:
	docker build -t $(APP_NAME):latest .

run: stop build
	docker compose up --build -d 

stop:
	docker compose down

docker-logs:
	docker compose logs -f
