---
name: opera-result-flow
description: "Use when generating or migrating Go error handling with go-opera: Result, Do, Try, Yield, MustPass, legacy if err checks, error wrapping, partial results, Map, FlatMap, and Catch recovery. Preserve legacy behavior rather than mechanically replacing every branch."
---

# Result Flow

Use package `opera` from `github.com/colafanta/go-opera` to compose fallible work
internally and return ordinary Go values and errors at public boundaries.

## Procedure

1. Read [migration patterns](./references/migration-patterns.md) for the legacy
   branch being replaced. Record its error-path value, error transformation,
   evaluation order, side effects, and deferred work.
2. Read [Result contracts](./references/result-api.md) before choosing helpers.
3. Use a small `Do(func() T)` block only for genuinely fail-fast work. Adapt
   `(T, error)` with `Try(...).Yield()` and error-only calls with
   `TryPass(...).Yield()`.
4. Preserve operation-specific wrapping at the operation; use block-level mapping
   or recovery only if it should apply to every matching error in the block.
5. End with `.Get()` or `.Err()` as appropriate. Do not silently discard partial
   error output or introduce a fallback for an operational failure.
6. Run focused tests for success, each early failure, the returned error-path
   value, `errors.Is`/`errors.As`, and calls that must not execute after failure.

Keep explicit branching for retries, accumulating errors, continuing after a
failure, or returning partial values. Do not treat fewer lines as equivalence.

## Usage Guidance

- Prefer functions that return `opera.Result[T]` for fallible application logic.
   Compose them inside `opera.Do` with `.Yield()`, and convert to ordinary Go
   `(T, error)` values at API boundaries with `.Get()`.
- When failure or absence should abort the operation, prefer
   `TryXX(...).Yield()` or `MayXX(...).Yield()` inside `opera.Do` over
   `MustXX(...)`. The explicit adapter and `Yield()` make the failure path easier
   to recognize, even when they take a few more lines. `Yield()` must only be
   called inside `opera.Do` so its bailout is converted to a `Result`.
- When a condition in the current `Do` block itself requires a deliberate
   bailout, use `opera.Err[opera.Unit](err).Yield()`. Reserve `opera.TryPass(err)`
   for adapting an operation whose only return value is `error`, such as
   `opera.TryPass(auditClient.Record(account.ID)).Yield()`.
- Use `opera.Try(value, err)` to adapt calls from other libraries into a
   `Do` block. When the target error is known, prefer specific helpers such as
   `MapErrIs`, `MapErrAs`, `CatchIs`, or `CatchAs` over general `MapErr` or
   `Catch`.

Assume the application types and clients below are defined elsewhere. Keep
application operations in the `Result` context, map known errors at their
source, and adapt ordinary library calls with `opera.Try`:

```go
func getUserAccount(id string) opera.Result[Account] {
   return opera.Try(accountClient.Find(id)).MapErrIs(sql.ErrNoRows, ErrAccountNotFound)
}

func getUserPreferences(account Account) opera.Result[Preferences] {
   return opera.Try(preferencesClient.Load(account.UserID))
}

func loadDashboard(id string) opera.Result[Dashboard] {
   return opera.Do(func() Dashboard {
      account := getUserAccount(id).Yield()
      if account.Disabled {
         opera.Err[opera.Unit](ErrAccountDisabled).Yield()
      }
      preferences := getUserPreferences(account).Yield()
      opera.TryPass(auditClient.Record(account.ID)).Yield()
      return opera.Try(dashboardClient.Build(account, preferences)).Yield()
   })
}

func handleDashboard(id string) (Dashboard, error) {
   return loadDashboard(id).Get()
}
```