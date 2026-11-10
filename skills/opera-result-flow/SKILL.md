---
name: opera-result-flow
description: "Use when generating or migrating Go error handling with go-opera: Result, Do, Try, TryYield, TryYieldUnit, YieldErr, legacy if err checks, error wrapping, typed error mapping and recovery, partial results, Map, FlatMap, and Catch. Preserve legacy behavior rather than mechanically replacing every branch."
---

# Result Flow

Use package `opera` from `github.com/colafanta/go-opera` to compose fallible work
internally and return ordinary Go values and errors at public boundaries.

## Procedure

1. Read [migration patterns](./references/migration-patterns.md) for the legacy
   branch being replaced. Record its error-path value, error transformation,
   evaluation order, side effects, and deferred work.
2. Read [Result contracts](./references/result-api.md) before choosing helpers.
3. Use a small `Do(func() T)` block only for genuinely fail-fast work. Prefer
   `TryYield(value, err)` for invoked `(T, error)` calls and
   `TryYieldUnit(operation())` for invoked error-only operations. Use
   `YieldErr[T](err)` when the error has
   already been constructed; in an early-return guard, explicitly return it.
   Use `Try(...).Yield()` or `TryPass(...).Yield()` when the Result needs
   composition or recovery before the bailout.
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
- When failure should abort the operation, prefer `opera.TryYield(value, err)`
   for ordinary `(T, error)` calls and `opera.TryYieldUnit(operation())` for
   invoked error-only calls inside `opera.Do`. Do not use deprecated `Must*` or
   `TryLookUp` / `TryAt` helpers in new code; use
   `MayLookUp(...).ToResult(err)` or `MayAt(...).ToResult(err)` for absence that
   must become an error. Use the `Result.Yield()` method only when a Result
   needs composition, mapping, or recovery first; the package-level
   `opera.Yield(result)` helper has been removed.
- For an already-constructed error, use `opera.YieldErr[T](err)`, not
   `opera.TryPass(err)` or another pass-through Unit wrapper. In an early-return
   guard, make the bailout explicit with `return` and the enclosing result type:

  ```go
  if err != nil {
      return opera.YieldErr[Dashboard](err)
  }
  ```

- Use `opera.TryPass(err)` when constructing a `Result[Unit]`, such as a workflow
   callback's return value. Do not use it just to bail out on a known error.
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
         return opera.YieldErr[Dashboard](ErrAccountDisabled)
      }
      preferences := getUserPreferences(account).Yield()
      opera.TryYieldUnit(auditClient.Record(account.ID))
      return opera.Try(dashboardClient.Build(account, preferences)).Yield()
   })
}

func handleDashboard(id string) (Dashboard, error) {
   return loadDashboard(id).Get()
}
```