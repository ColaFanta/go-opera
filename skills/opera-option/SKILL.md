---
name: opera-option
description: "Use when generating or migrating go-opera Option code for expected absence, map lookups, slice bounds, nullable pointers, empty values, JSON omission, SQL Scan/Value, or serialization. Distinguish missing data from operational errors and preserve valid zero values."
---

# Optional Values

`Option[T]` models presence without an error, not a failed operation.

## Procedure

1. Decide whether absence is expected, an error, or distinct from a valid zero
   value. Read [presence patterns](./references/presence-patterns.md).
2. Choose a constructor that preserves that distinction. Do not use
   `MaybeEmpty` when `0`, `false`, or an empty string is valid present data.
3. Keep absence optional with `Get`, `Or`, or `OrElse`. If it must fail, use
   `ToResult()` with call-local error mapping and then `Yield` inside `Do`.
4. For nullable fields, read [integration contracts](./references/integrations.md)
   before changing tags, database bindings, or serialized representation.
5. Test present nonzero, present zero, missing, and malformed input separately.
   For migrations, preserve error identity, fallback evaluation, and pointer
   aliasing requirements.

`None` and a zero `Option[T]` have no payload. `Some(zero)` is present and must
not be mistaken for absence by application logic.