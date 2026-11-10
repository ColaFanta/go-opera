# Serialization And Database Contracts

Option supports JSON, text, binary, gob, SQL `Scan`, and driver `Value` adapters.
Check the concrete payload's encoder/decoder and database compatibility before
replacing a pointer, `sql.Null*`, or an existing custom nullable type.

## JSON

```go
type Update struct {
    Name opera.Option[string] `json:"name,omitzero"`
}
```

- Some marshals its contained value; None marshals as JSON null when not omitted.
- Use `omitzero` on Go 1.24+ rather than assuming `omitempty` omits a custom
  marshaling struct. Option implements `MarshalJSON` and `IsZero`.
- `IsZero` can also be true for a present zero-valued payload such as
  `Some("")`. Consequently `omitzero` does not mean "omit None only".
- An omitted field and an explicit null are not automatically a three-state
  patch model. If the API distinguishes unchanged, cleared, and assigned,
  preserve that distinction with explicit request modeling.
- Test missing, null, present zero, present nonzero, invalid input, and reuse of
  an existing destination before claiming wire-format equivalence.

## SQL And Other Encodings

- `Scan` and `Value` provide database integration, but an Option is not a blanket
  replacement for every nullable driver type or custom scanner.
- Verify NULL behavior, supported source/driver types, delegated payload methods,
  and errors using the actual database binding or codec.
- Do not turn a decode or query error into None with `ToOption` unless losing
  that failure is the existing contract. Operational errors belong in Result.
- Preserve round-trip expectations separately from presence checks: Some(zero)
  and None may encode or omit differently depending on tags and codec.