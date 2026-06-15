FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /transfer-system .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /transfer-system /app/transfer-system

ENV DATABASE_URL=postgres://transfer_system:transfer_system@postgres:5432/transfer_system?sslmode=disable

EXPOSE 8080

CMD ["/app/transfer-system"]
