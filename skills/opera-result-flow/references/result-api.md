# Result Contracts

`Result[T]` stores a value and an error. Success means `err == nil`; errors may
carry nonzero partial values. The zero Result is a successful zero value.

## Constructors And Bailouts

| Adapter | Result | Bailout Form Inside `Do` |
| --- | --- | --- |
| `Try(value, err)` | `Result[T]`, retains value even on error | `Must(value, err)` |
| `TryPass(err)` | `Result[Unit]` | `MustPass(err)` |
| `TryHave(value, ok)` | Missing wraps `ErrNoSuchElement`; retains value | `MustHave(value, ok)` |
| `TryLookUp(m, key)` | Missing wraps `ErrNoSuchElement` | `MustLookUp(m, key)` |
| `TryAt(slice, index)` | Negative/out-of-range index wraps `ErrNoSuchElement` | `MustAt(slice, index)` |
| `TryIf(ok)` | False gives `ErrPredicateFalse` | `MustTrue(ok)` |
| `TryCast[T](value)` | Assertion failure wraps `ErrInputCastFailed` | `MustCast[T](value)` |

`Ok(value)` constructs success. `Err[T](err)` constructs a zero-valued payload
with the supplied error; `Err[T](nil)` is therefore successful, not a failure.
Use `Empty[T]()` for T's zero value and `opera.U` for a `Unit` payload.

## Inspection And Extraction

- `IsOk()`, `IsErr()`: inspect error presence.
- `Get() (T, error)`: exposes both fields, including error-path output.
- `Err() error`: returns only the error.
- `Yield() T`: returns success or panics with the internal bailout wrapper.
- `Or(fallback T)`, `OrF(func(error) T)`, `OrEmpty()`: use a fallback on failure,
  never the partial error payload. `OrF` is lazy; `Or` arguments are eager.
- `OrPanic() T`: panics with the original error on failure, not the bailout
  wrapper. `Do` does not convert this panic into a Result.
- `ToPointer()`, `ToOption()`: convert success; errors become nil/None and lose
  the error and any partial value. Do not use when failure must be reported.

`Do(func() T) Result[T]` catches only the internal wrapper. A caught bailout
returns the block output type's zero value and the error, not the failed
operation's partial value. `Option.Yield` uses the same wrapper for absence.

## Composition

- `Map[M](func(T) M) Result[M]`: transform success.
- `FlatMap[M](func(T) Result[M]) Result[M]`: chain success.
- `Tap(func(T))`: success-only side effect; returns the original Result.
- `TapErr(func(error))`: error-only side effect; returns the original Result.
- `Cast[C]() Result[C]`: assert the value to another type. An existing error
  retains its payload only when assertable to C. A failed assertion on success
  returns an error wrapping `ErrResultCastFailed`.

This checkout supports type-changing generic methods. On an existing error,
`Map` and `FlatMap` skip their callback and retain the payload only if it can be
asserted to the destination type; otherwise the destination value is zero.
Read partial output with `Get()` before changing types if it must be retained.

## Error Mapping And Recovery

- `MapErr(func(error) error)`: replace an existing error, retaining the value.
- `MapErrIs(target, replacement)`: map an `errors.Is` match.
- `MapErrAs[E error](func(E) error)`: map a matching typed error using
  `errors.AsType`.
- `MapErrString(substring, replacement)`: map a message substring match.
- `Catch(func(error) Result[T])`: lazy recovery on error.
- `CatchAs[E error](func(E) Result[T])`: lazy recovery on the first typed match
  found by `errors.AsType`, including wrapped and joined errors. The handler
  receives the matched error, not the outer wrapper; E can be inferred from its
  parameter type. Successes and unmatched failures are unchanged.
- `CatchIs(target, fallback)`: replace a matching failure with successful fallback.
- `CatchString(substring, fallback)`: recover a message substring match.

Prefer sentinel matching over message matching when the upstream API supports it.
`ErrAny` intentionally matches any existing error in the matching helpers.
Replacing an error with nil changes failure into success; make that explicit.

If a `Catch` or `CatchAs` handler returns another failure, it retains the original
Result's payload with the handler's error, not the handler's failure payload.

Do not invent `Unwrap`, `UnwrapOr`, `UnwrapOrElse`, `AndThen`, `Expect`, or
`Recover`. Those are not this package's API.