# Behavior-Preserving Migration Patterns

Examples assume an import of `github.com/colafanta/go-opera` as `opera` and the
application's existing types and dependencies. Keep public function signatures.

## Fallible Value Calls

Before:

```go
func Load(id string) (Item, error) {
    row, err := repo.Find(id)
    if err != nil {
        return Item{}, err
    }
    return build(row), nil
}
```

After:

```go
func Load(id string) (Item, error) {
    return opera.Do(func() Item {
        row := opera.Try(repo.Find(id)).Yield()
        return build(row)
    }).Get()
}
```

Equivalent only if the legacy error-path output is the output type's zero value,
the error is propagated unchanged, and the same calls run in the same order.
`opera.Must(repo.Find(id))` is the shorter equivalent inside `Do`.

## Error-Only Sequences

Before:

```go
func Save() error {
    if err := validate(); err != nil {
        return err
    }
    if err := persist(); err != nil {
        return err
    }
    return nil
}
```

After:

```go
func Save() error {
    return opera.Do(func() opera.Unit {
        opera.MustPass(validate())
        opera.MustPass(persist())
        return opera.U
    }).Err()
}
```

This preserves stop-on-first-error behavior. It does not replace a loop that
collects errors or intentionally continues processing.

## Preserve Wrapping And Match Scope

To replace `return Item{}, fmt.Errorf("load %s: %w", id, err)`, wrap at the call:

```go
row := opera.Try(repo.Find(id)).MapErr(func(err error) error {
    return fmt.Errorf("load %s: %w", id, err)
}).Yield()
```

Moving this wrapper after the entire `Do` also wraps failures from later calls,
which is not equivalent. Preserve `%w`, error identity, message context, and
whether the legacy code used direct equality or `errors.Is`.

## Known Benign Errors

For an original `errors.Is(err, sql.ErrNoRows)` fallback on a single call:

```go
func FindOrDefault(id string) (Item, error) {
    return opera.Try(repo.Find(id)).CatchIs(sql.ErrNoRows, defaultItem).Get()
}
```

The fallback expression is evaluated before `CatchIs` runs. If constructing it
has side effects, is expensive, or can fail, retain a branch or use lazy `Catch`:

```go
return opera.Try(repo.Find(id)).Catch(func(err error) opera.Result[Item] {
    if errors.Is(err, sql.ErrNoRows) {
        return opera.Try(loadDefault())
    }
    return opera.Err[Item](err)
}).Get()
```

A `CatchIs` after a multi-operation `Do` applies to every matching error from the
block. Use it there only when that broader recovery is intended. Never use
`Or` or `Catch` merely to hide an unhandled operational error.

## Partial Output Is Not Fail-Fast Output

These return patterns must not be mechanically converted to `Yield`:

```go
data, err := readPartial()
if err != nil {
    return data, err
}
return data, nil
```

`Try` preserves both fields, so `return opera.Try(readPartial()).Get()` is
equivalent. In contrast, `Do(func() []byte { return
opera.Try(readPartial()).Yield() }).Get()` loses the partial bytes on failure.
The same risk applies to nonzero sentinel outputs, counts, accumulated state,
and `io.Reader` results where `n > 0` and `err != nil` may coexist.

Use `.Get()` and explicit branching when both the partial output and error are
needed. Type-changing transforms may also discard error payloads.

## Presence And Predicates

| Legacy Pattern | Candidate | Equivalence Condition |
| --- | --- | --- |
| `(value, ok)` with failure on `!ok` | `TryHave(value, ok).Yield()` | Its wrapped `ErrNoSuchElement` matches the required error contract |
| Map lookup with required key | `TryLookUp(m, key).Yield()` or `MustLookUp(m, key)` | Preserve custom errors with call-local mapping |
| Bounds-checked slice access | `TryAt(items, index).Yield()` or `MustAt(items, index)` | Absence was an error, not a panic or a default |
| Required boolean condition | `TryIf(ok).Yield()` or `MustTrue(ok)` | `ErrPredicateFalse` matches or is explicitly remapped |
| Checked type assertion | `TryCast[T](value).Yield()` or `MustCast[T](value)` | Preserve prior failure behavior and error contract |

For custom errors, map at that exact operation, for example:

```go
item := opera.TryLookUp(items, key).MapErrIs(opera.ErrNoSuchElement, ErrMissingItem).Yield()
opera.TryIf(allowed).MapErrIs(opera.ErrPredicateFalse, ErrForbidden).Yield()
```

Do not replace an unchecked assertion or indexing panic with an error unless
changing the contract is intentional. See the optional-values skill for expected
absence rather than failure.

## Control Flow And Defers

- Keep guards, branching, loops, retry decisions, and error accumulation when
  those express the behavior. `Do` is not a general exception handler.
- A defer moved inside a `Do` callback runs when that callback exits, not when
  the surrounding function returns. Preserve resource lifetime and side effects.
- Deferred functions that modify named returns do not automatically modify the
  enclosing `Result`; retain the original structure when necessary.
- `Do` only catches opera's internal bailout wrapper. `panic(err)` and
  `panic("bug")` are re-panicked. Do not use it to replace a general `recover`.
- Test failure at every migrated operation, not just the successful pipeline.