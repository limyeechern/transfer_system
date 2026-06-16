# Transfer System

## Assumptions

### Transfer Amount

Transfer amounts are to be stored as fixed-scale integers using `BIGINT`, normalized to 5 decimal places in the application layer.

This means that system assumes:

- Transfers do not exceed approximately `92 trillion`.
- Amount precision does not exceed 5 decimal places.

For example, `123.45678` is stored as `12345678`.

When amounts are returned by the API, they are formatted back into decimal strings with exactly 5 decimal places. For example, a stored value of `12345678` is returned as `"123.45678"`, and `10000000` is returned as `"100.00000"`.

### Transaction Idempotency

Transaction creation is currently not idempotent. The API does not accept a client-provided unique request id or idempotency key, so retries may create a new transaction.

### Read/Write Skew

The system assumes account reads are significantly more frequent than writes. Redis is used as a read-through cache for account lookups to reduce repeated database reads while Postgres remains the source of truth.
