# Presence Patterns And API

## Choosing A Constructor

| Existing Input | Adapter | Important Constraint |
| --- | --- | --- |
| Known present value | `Some(value)` | Zero values remain present |
| Known missing value | `None[T]()` | Canonical zero payload, not retained data |
| `(value, ok)` | `MayHave(value, ok)` | False discards the value |
| Map lookup | `MayLookUp(m, key)` | A present zero value remains present |
| Checked slice index | `MayAt(items, index)` | Negative/out-of-bounds returns None |
| Nullable `*T` | `MaybeNilPtr(ptr)` | Dereferences into an Option value, not pointer identity |
| Zero means missing | `MaybeEmpty(value)` | Uses broad `IsEmpty`, not just empty-string checking |
| Checked assertion | `MayCast[T](value)` | Failed assertion becomes None, not an error |

`IsEmpty` handles zero values and type-specific zero checks. Do not replace a
narrow predicate such as `s == ""` or `len(items) == 0` with `MaybeEmpty` without
checking equivalence for the actual type.

## Optional Lookup With Fallback

Before:

```go
name, ok := names[id]
if !ok {
    name = "anonymous"
}
```

After:

```go
name := opera.MayLookUp(names, id).Or("anonymous")
```

An existing empty string remains present. `MaybeEmpty(names[id])` would not
preserve that behavior. `Or(fallback())` evaluates `fallback()` eagerly; retain
the branch or use lazy `OrElse` if that changes side effects:

```go
name := opera.MayLookUp(names, id).OrElse(func() opera.Option[string] {
    return opera.Some(makeFallbackName())
}).Yield()
```

The last form uses `Yield`; put it inside `Do`, or use `Get`/`OrEmpty` instead
when no managed bailout context exists.

## Required Missing Value With A Domain Error

Before:

```go
name, ok := names[id]
if !ok {
    return User{}, ErrNameMissing
}
```

After, inside `Do(func() User)`:

```go
name := opera.MayLookUp(names, id).ToResult().
    MapErrIs(opera.ErrNoSuchElement, ErrNameMissing).Yield()
```

Direct `Option.Yield()` also aborts on absence, but gives `ErrNoSuchElement`, not
the legacy domain error. `ToResult` wraps that sentinel with type information;
use `errors.Is` when matching it. Call-local mapping avoids accidentally
rewriting an unrelated missing value elsewhere in the block.

## Nullable Pointer

`opera.MaybeNilPtr(input.Name).Or("anonymous")` replaces a nil check plus
dereference when a value copy is intended. It does not retain the original
pointer's identity or make later mutations affect the input.

For a required pointer, use `MaybeNilPtr(ptr).ToResult()` and map absence to the
original error inside `Do`. Preserve explicit pointer handling when aliasing,
mutation, or object ownership matters.

## Option Methods

- Inspection: `IsSome`, `IsNone`, `Get() (T, bool)`.
- Extraction: `Yield`, `Or`, `OrEmpty`, `OrPanic`.
- Composition: `Map[M](func(T) M)`, `FlatMap[M](func(T) Option[M])`.
- Side effects: `Tap(func(T))`, `TapNone(func())`.
- Lazy alternative: `OrElse(func() Option[T])`.
- Conversion: `ToResult`, `ToPointer`, `Cast[C]`.
- Integration: `IsZero` supports zero-value omission, not presence testing.

`Map` and `FlatMap` may change type. On None they skip the callback and return
canonical `None[M]()`. `Get` on None returns the type's zero value and false;
`ToPointer` returns nil. `ToPointer` on Some returns a pointer to a copy of the
stored value, not an original input pointer.

`Yield` on None uses the bailout wrapper; `OrPanic` instead panics directly with
`ErrNoSuchElement`, which `Do` does not convert into a Result. There is no
`Option.Equal` method in this checkout.