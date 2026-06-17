.PHONY: run build test test-integration postgres-test fmt clean docker-build docker-up docker-down docker-logs

APP_NAME := transfer-system
TEST_DATABASE_URL := postgres://transfer_system:transfer_system@127.0.0.1:15433/transfer_system_test?sslmode=disable

test:
	go test ./...

postgres-test:
	docker compose up -d --wait postgres-test

test-integration: postgres-test
	TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test ./biz/dal -v
	docker compose stop postgres-test

fmt:
	go fmt ./...

clean:
	rm -rf bin

build:
	docker build -t $(APP_NAME):latest .

run: stop build
	docker compose up --build -d 

stop:
	docker compose down -v

docker-logs:
	docker compose logs -f
