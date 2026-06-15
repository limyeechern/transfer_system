# Transfer System

## Assumption

### Transfer Amount

Transfer amounts are assumed to be stored as fixed-scale integers using `BIGINT`, normalized to 5 decimal places in the application layer.

This means the system assumes:

- Transfers do not exceed approximately `92 trillion`.
- Amount precision does not exceed 5 decimal places.

For example, `123.45678` is stored as `12345678`.

When amounts are returned by the API, they are formatted back into decimal strings with exactly 5 decimal places. For example, a stored value of `12345678` is returned as `"123.45678"`, and `10000000` is returned as `"100.00000"`.
