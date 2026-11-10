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
   `(T, error)` with `Try(...).Yield()` and error-only calls with `MustPass(...)`.
4. Preserve operation-specific wrapping at the operation; use block-level mapping
   or recovery only if it should apply to every matching error in the block.
5. End with `.Get()` or `.Err()` as appropriate. Do not silently discard partial
   error output or introduce a fallback for an operational failure.
6. Run focused tests for success, each early failure, the returned error-path
   value, `errors.Is`/`errors.As`, and calls that must not execute after failure.

Keep explicit branching for retries, accumulating errors, continuing after a
failure, or returning partial values. Do not treat fewer lines as equivalence.