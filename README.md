# Transfer System

## Assumption

### Transfer Amount

Transfer amounts are assumed to be stored as fixed-scale integers using `BIGINT`, normalized to 5 decimal places in the application layer.

This means the system assumes:

- Transfers do not exceed approximately `92 trillion`.
- Amount precision does not exceed 5 decimal places.

For example, `123.45678` is stored as `12345678`.
