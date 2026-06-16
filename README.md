# Transfer System

## Instructions

### Installation

Install:

- Go `1.20`
- Docker
- Docker Compose

Download Go dependencies:

```sh
go mod download
```

### Set Up

```sh
cp example.env .env
```

The default local services are:

- API: `http://127.0.0.1:8080`
- Postgres: `127.0.0.1:15432`
- Redis: `127.0.0.1:6379`

Default environment variables:

```sh
DATABASE_URL=postgres://transfer_system:transfer_system@127.0.0.1:15432/transfer_system?sslmode=disable
REDIS_ADDR=127.0.0.1:6379
```

Docker Compose starts Postgres and Redis for local development.

### Run

Start the system:

```sh
make run
```

View logs:

```sh
make docker-logs
```

Run tests:

```sh
go test ./...
```

### API Examples

Create an account:

```sh
curl -X POST http://127.0.0.1:8080/accounts \
  -H 'Content-Type: application/json' \
  -d '{"account_id":123,"initial_balance":"100.23344"}'
```

Get an account:

```sh
curl http://127.0.0.1:8080/accounts/123
```

Create a transaction:

```sh
curl -X POST http://127.0.0.1:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{"source_account_id":123,"destination_account_id":456,"amount":"10.12345"}'
```

## Assumptions

### 1. Transfer Amount

Transfer amounts are to be stored as fixed-scale integers using `BIGINT`, normalized to 5 decimal places in the application layer.

This means that system assumes:

- Transfers do not exceed approximately `92 trillion`.
- Amount precision does not exceed 5 decimal places.

For example, `123.45678` is stored as `12345678`.

When amounts are returned by the API, they are formatted back into decimal strings with exactly 5 decimal places. For example, a stored value of `12345678` is returned as `"123.45678"`, and `10000000` is returned as `"100.00000"`.

### 2. Transaction Idempotency

Transaction creation is currently not idempotent. The API does not accept a client-provided unique request id or idempotency key, so retries may create a new transaction.

### 3. Read/Write Skew

The system assumes account reads are significantly more frequent than writes. Redis is used as a read-through cache for account lookups to reduce repeated database reads while Postgres remains the source of truth.

### 4. Concurrency

The system assumes that the currency is the same for all accounts.

### 5. Authentication

No need to implement authn or authz
