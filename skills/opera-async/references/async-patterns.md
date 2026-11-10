# Async Contracts And Patterns

## Independent Tasks

```go
func LoadBoth(ctx context.Context) (opera.Tuple[int, string], error) {
    return opera.Do(func() opera.Tuple[int, string] {
        countTask := opera.Async(ctx, func(ctx context.Context) opera.Result[int] {
            return opera.Try(loadCount(ctx))
        })
        nameTask := opera.Async(ctx, func(ctx context.Context) opera.Result[string] {
            return opera.Try(loadName(ctx))
        })
        return opera.AwaitZip(ctx, countTask, nameTask).Yield()
    }).Get()
}
```

Both tasks start before waiting. This is equivalent to legacy concurrent work
only when its error selection, lifetime, and cancellation rules match. Changing
sequential calls into this form changes side-effect order and may start the
second operation even when the first fails. Treat that as a deliberate behavior
change, not cleanup.

## Choosing A Wait

| Helper | Contract |
| --- | --- |
| `Await(ctx, ch)` | One Result or wait-context cancellation |
| `AwaitAllT(ctx, chans...)` | Homogeneous values in input order; returns early on error |
| `AwaitAllSettledT(ctx, chans...)` | Result per input, including failures; no failure short-circuit |
| `AwaitFirstT(ctx, chans...)` | First successful value and original channel index |
| `AwaitZip`, `AwaitZip3` ... `AwaitZip9` | Typed tuple from differently typed channels |
| `AwaitAll[any]`, `AwaitFirst[any]`, `AwaitAllSettled[any]` | Dynamic heterogeneous channels when a typed tuple is unsuitable |

Dynamic `AwaitAll[T]` and `AwaitFirst[T]` accept `...any`; incompatible payload
types can produce `ErrAsyncAwaitedValueCastError`. Prefer typed wrappers for
homogeneous inputs. Zip helpers cover two through nine channels, not `AwaitZip1`.

`AwaitFirst` returns `Result[Tuple[T, int]]`. It ignores failed candidates while
success remains possible. If all active candidates fail, the joined error
includes `ErrAsyncAllOpsFailed`. Simultaneously ready successes may win
nondeterministically; do not depend on tie ordering.

## Nil And Closed Channels

- `Await` on a nil channel returns `Ok(zero T)` immediately, unlike a raw Go
  receive from nil, which blocks forever.
- `AwaitAll` and `AwaitAllSettled` treat nil inputs as successful zero entries.
- `AwaitFirst` ignores nil candidates; with all inputs nil it succeeds with a
  zero tuple, not a missing-task error.
- `Await` on a closed channel without a value returns an error wrapping
  `ErrAsyncChannelClosed`, not a successful zero value.
- Validate missing required tasks explicitly rather than relying on nil to fail.

## Cancellation And Ownership

- `Async(ctx, thunk)` accepts `Thunk[T]`, a function from `context.Context` to
  `Result[T]`, and returns a receive-only, buffered, one-result channel.
- The thunk executes in a goroutine and must cooperate with its context itself.
  Cancellation is not a guarantee that the thunk will never start or will stop.
- The implementation selects between result delivery and cancellation; it does
  not guarantee cancellation always wins when both are ready.
- Async does not close its result channel. Await a task once; do not range over
  it expecting closure or assume the channel broadcasts to multiple consumers.
- Await-context cancellation ends waiting with an error. Fail-fast joins and
  first-success waits do not automatically cancel or join sibling producers.
- If this operation owns the producers, derive a context with
  `context.WithCancel`, pass it to those producers, and call its cancel function
  at the correct lifetime boundary. Retain explicit joining if all producers
  must exit before returning; cancellation alone is not joining.
- Settled waits collect outcomes rather than short-circuiting on task failures;
  a canceled wait can still yield cancellation errors instead of completed
  producer results. Do not mistake it for a guarantee producers have exited.
- Async is not a panic-to-error boundary. Wrap fallible linear code in `Do` in
  the thunk if needed; ordinary panics still propagate.

Do not replace streaming channels, broadcasts, worker pools, or bounded
concurrency controls with these one-result helpers without retaining those
protocols and limits.

## Adapting Parameterized Functions

`Suspend0`, `Suspend`, and `Suspend2` through `Suspend9` adapt functions taking a
context and zero through nine application arguments into factories of `Thunk`.
They expect Result-returning functions; adapt ordinary `(T, error)` functions
with `Try` first. Use a closure when that is clearer than an adapter.