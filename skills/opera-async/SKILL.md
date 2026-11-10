---
name: opera-async
description: "Use when generating or migrating go-opera concurrent tasks: Async, Await, AwaitAllT, AwaitAllSettledT, AwaitFirstT, AwaitZip, Suspend, context cancellation, channel ownership, and ordered joins. Check concurrency equivalence before replacing goroutines or sequential code."
---

# Async Coordination

Use async helpers for independent Result-returning work, not as a mechanical
replacement for sequential calls or every goroutine/channel protocol.

## Procedure

1. Identify dependencies, shared state, call ordering, cancellation ownership,
   and whether every task must finish before the function returns.
2. Read [async contracts and patterns](./references/async-patterns.md).
3. Choose `Await` for one task, `AwaitAllT` for homogeneous fail-fast joining,
   `AwaitAllSettledT` for every outcome, `AwaitFirstT` for first success, or
   `AwaitZip`/`AwaitZip3` through `AwaitZip9` for heterogeneous typed tuples.
4. Pass context to the actual work. Keep cancellation of sibling producers
   explicit; stopping a wait does not stop or join those producers.
5. Extract through `Yield` only inside a managed `Do`, and return `.Get()` or
   `.Err()` at the conventional boundary.
6. Test errors, cancellation, nil/closed channels, ordering, and any required
   completion guarantee. Use the race detector when shared mutable state changes.