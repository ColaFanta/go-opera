# opera

![Go Version](https://img.shields.io/github/go-mod/go-version/ColaFanta/go-opera)
![Test Status](https://github.com/ColaFanta/go-opera/actions/workflows/unit-test.yml/badge.svg)


`opera` - **Op**tion, **E**rror handling, **R**esult, and **A**sync, a library that brings better error handling and async programming to Golang.

## Acknowledgments

This library was inspired by two excellent projects:

- [eh](https://github.com/olevski/eh), which introduced the idea of Rust‑like error handling in Go.
- [mo](https://github.com/samber/mo), from which the `Do` notation and `Option` marshaling into/out of values were adapted.

Many thanks to their authors for the inspiration.

## Use Case
### Error Handling with opera
Checking a username and password often involves nested if statements and explicit error handling, for example:

```go
func checkLogin(user, password string) (bool, error) {
     ctx, _ := context.Background()
 
    // unpack result in the old way
    u, err := gorm.G[model.User](db).
        Where(gen.User.Username.Eq(user)).
        First(ctx)
    if err != nil {
        // return false if user not found
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
        // propagate other errors
        return false, err
    }

     if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
        return false, nil
    }

     return true, nil
}
```
With opera `Do` notation, you can safely unwrap the result of a failable function without repeatedly writing `if val, err := fn(); err != nil`. The output of a function can be safely produced with `Result.Yield()`, and fast‑fail is implicitly handled by `Do` notation. No explicit manipulation on `error` is needed in the context of `Do` notation.
:

```go
func(user, password string) (ok bool, err error) {
    // use `Do` notation
    return opera.Do(func() bool {
        ctx, _ := context.Background()

        // get user with `Do` notation and `Result` without `if`
        // fail fast to `Do` when `err` occurs
        user := opera.Try(gorm.G[model.User](db).
            Where(gen.User.Username.Eq(user)).
            First(ctx)).Yield()

        // fail fast also applies to `Option` type
        pwd := opera.MaybeNilPtr(user.Password).Yield()

        // also fail fast if seen `ErrMismatchedHashAndPassword`
        opera.Must0(bcrypt.CompareHashAndPassword([]byte(pwd), []byte(password)))

        // login success
        return true
    }).
        // conditionally recover if certain errors are met 
        CatchIs(gorm.ErrRecordNotFound, false).
        CatchIs(bcrypt.ErrMismatchedHashAndPassword, false).
        Get() // unpack to `(T, error)` in the end
}
```

#### Performance Caveat
The design of how this library achieves fail‑fast error handling may introduce little performance overhead. For reference, below is a simple benchmark result of two test cases in the project comparing the same operation with and without `Do`.

```
cpu: AMD Ryzen 9 8945HS w/ Radeon 780M Graphics
BenchmarkDoAndGet/ErrFromGet-16              6267090           188.3 ns/op         128 B/op           3 allocs/op
BenchmarkDoAndGet/ErrFromDo-16               6167328           196.9 ns/op         128 B/op           3 allocs/op
```

### Concurrency with opera `Async`
`Async` in opera provides a clean way to organize structured concurrency tasks without manually managing goroutines, channels, or synchronization primitives.

```go
func(ctx context.Context, q *gorm.DB, p Params) error {
    return opera.Do(func () any {
        // ctx `Done` of parent also cancels tasks fired by opera.Async and unblocks opera.Await
        countTask := opera.Async(ctx, func(ctx context.Context) opera.Result[int64] {
            return opera.Try(q.Count(ctx, "*"))
        })
        dataTask := opera.Async(ctx, func(ctx context.Context) opera.Result[[]T] {
            size := p.Size.Or(20)
            offset := (p.Page.Or(1) - 1) * size
            q.Offset(offset).Limit(size)
            return opera.Try(q.Find(ctx))
        })
        // Getting results of two async tasks in one line
        count, data := opera.Await(ctx, countTask).Yield(), opera.Await(ctx, dataTask).Yield()
        return nil
    }).Err()
}
```