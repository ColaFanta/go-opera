# Cleanup Equivalence And Patterns

## Scoped

Signature:

```go
func Scoped[T, R any](
    acquire func() opera.Result[T],
    use func(T) opera.Result[R],
    release func(T),
) opera.Result[R]
```

This signature describes the exported function, not code to redeclare locally.
`Scoped` gets the acquisition result, defers `release(acquiredValue)`, then checks
the acquisition error. Release runs even if acquisition failed, receiving its
partial value or zero value. Use must not run after acquisition failure.

The common legacy shape below is NOT directly equivalent:

```go
resource, err := acquire()
if err != nil {
    return Output{}, err
}
defer release(resource)
return use(resource)
```

Here release is registered only on success. Use `Scoped` only when release on
failed acquisition is intended or is demonstrably harmless, including counters,
logs, nil handling, and partial resource ownership. Its `release func(T)` cannot
return an error; do not discard an important close/commit error to fit that API.

Release runs when `Scoped` exits, not at the end of the outer function. It runs
on a panic from use because it is deferred. Do not confuse these semantics with
workflow finalizers.

## Workflow

`Workflow()` returns an unexported `*workflow`. Infer its type in a local
variable; do not name `opera.workflow` or invent an exported workflow type.

`Step(flow, work, finalizer)`:

1. Registers `func(error) Result[Unit]` BEFORE calling work.
2. Executes `work func() Result[T]`.
3. If work returns an error, runs all registered finalizers in reverse order,
   including the failing step's finalizer, joining finalizer errors.
4. Yields the work result, aborting a surrounding `Do` on failure.

`Finalize(flow)` runs registered finalizers in reverse order on the success
path. It starts with a nil error; if a finalizer fails, later finalizers receive
the accumulated error. Cleanup that branches on `err != nil` can therefore run
rollback even during success-path finalization after an earlier cleanup failure.
Finalization errors bail out through `MustPass`.

Do not reuse a flow after finalization or retry failed finalization. The package
is not a retryable transaction manager.

## Safe Closure Scope

Declare captured outputs before passing closures to `Step`. A variable declared
by `id := Step(... closure using id ...)` is not in scope in that initializer.

```go
func Provision() error {
    return opera.Do(func() opera.Unit {
        flow := opera.Workflow()
        var userID string
        var userCreated bool

        opera.Step(flow,
            func() opera.Result[string] {
                return opera.Try(createUser()).Tap(func(id string) {
                    userID = id
                    userCreated = true
                })
            },
            func(err error) opera.Result[opera.Unit] {
                if err != nil && userCreated {
                    return opera.TryPass(deleteUser(userID))
                }
                return opera.Ok(opera.U)
            },
        )

        opera.Step(flow,
            func() opera.Result[opera.Unit] {
                return opera.TryPass(configureUser(userID))
            },
            func(error) opera.Result[opera.Unit] {
                return opera.Ok(opera.U)
            },
        )

        opera.Finalize(flow)
        return opera.U
    }).Err()
}
```

This assumes failed `createUser` calls need no cleanup of partial output and
`configureUser` needs no separate compensation. Track partial resources with
`Get()` and add real compensators when those assumptions do not hold.

## When Not To Replace Legacy Cleanup

- A panic in work does not run workflow finalizers as deferred Go cleanup would.
- An unrelated `MustPass` or `Yield` between successful steps can abort the
  enclosing `Do` without triggering workflow finalizers. Put fallible operations
  that need this rollback inside `Step`.
- A failing step's finalizer runs even if work acquired nothing. Preserve
  acquisition flags or ownership state rather than deleting an unowned resource.
- `errors.Join` changes cleanup-error reporting compared with ignoring cleanup
  errors or returning only the original failure. Preserve or explicitly approve
  that contract change.
- Finalizers receive errors, not the step output. Capture state safely and do
  not assume a successful output exists after a failed step.
- Keep native transaction APIs for atomic database commits and rollbacks.
  Compensation across external systems is not atomicity.