# Result Contracts

`Result[T]` stores a value and an error. Success means `err == nil`; errors may
carry nonzero partial values. The zero Result is a successful zero value.

## Constructors And Bailouts

| Adapter | Result | Bailout Form Inside `Do` |
| --- | --- | --- |
| `Try(value, err)` | `Result[T]`, retains value even on error | `TryYield(value, err)` |
| `TryPass(err)` | `Result[Unit]` | `TryYieldUnit(operation())` for an invoked error-only operation; `YieldErr[T](err)` for an already-constructed error |
| `TryHave(value, ok)` | Missing wraps `ErrNoSuchElement`; retains value | `TryHave(value, ok).Yield()` |
| `MayLookUp(m, key).ToResult(err)` | Option absent becomes the supplied Result error | `.Yield()` in `Do` |
| `MayAt(slice, index).ToResult(err)` | Option absent becomes the supplied Result error | `.Yield()` in `Do` |
| Direct condition check | Use `if` and `Err[Unit](err).Yield()` | No legacy `TryIf`/`MustTrue` adapter |
| `TryCast[T](value)` | Assertion failure wraps `ErrInputCastFailed` | `TryCast[T](value).Yield()` |

`TryLookUp`, `TryAt`, `TryIf`, and every `Must*` function are deprecated. Do not
introduce them in new code. `TryYield` is a direct shortcut for
`Try(value, err).Yield()`; `TryYieldUnit` is the equivalent for an invoked
error-only operation. For an error value already in hand, use `YieldErr[T](err)`;
do not wrap it in `TryPass` or another pass-through Unit wrapper merely to
trigger a bailout. In an early-return guard, explicitly `return YieldErr[T](err)`
using the enclosing output type.

`Ok(value)` constructs success. `Err[T](err)` constructs a zero-valued payload
with the supplied error; `Err[T](nil)` is therefore successful, not a failure.
Use `Empty[T]()` for T's zero value and `opera.U` for a `Unit` payload.

## Inspection And Extraction

- `IsOk()`, `IsErr()`: inspect error presence.
- `Get() (T, error)`: exposes both fields, including error-path output.
- `Err() error`: returns only the error.
- `Try2` through `Try9`: wrap 2-9 ordinary `(values, err)` results into tuples.
- `TryYield(value, err)`: inside `Do`, adapt and yield a value/error call directly.
- `TryYieldUnit(err)`: inside `Do`, adapt and yield an error-only call directly.
- `YieldErr[T](err)`: inside `Do`, bail out with an already-constructed error;
  `nil` is misuse and panics with `ErrNilYieldError`.
- `Yield() T`: returns success or panics with the internal bailout wrapper.
- `Or(fallback T)`, `OrF(func(error) T)`, `OrEmpty()`: use a fallback on failure,
  never the partial error payload. `OrF` is lazy; `Or` arguments are eager.
- `OrPanic() T`: panics with the original error on failure, not the bailout
  wrapper. `Do` does not convert this panic into a Result.
- `ToPointer()`, `ToOption()`: convert success; errors become nil/None and lose
  the error and any partial value. Do not use when failure must be reported.

There is no package-level `opera.Yield(result)` helper; call the
`Result.Yield()` method when a composed Result must bail out.

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