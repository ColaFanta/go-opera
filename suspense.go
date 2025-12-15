package opera

import "context"

// Thunk is just a type alias for any function that defers computation until invoked with a context.
// A `Thunk` can yield `Result` with the help of `opera.Async`.
//
// Example:
//
//	thunk := func(ctx context.Context) opera.Result[int] {
//		// perform some computation
//		return opera.Ok(42)
//	}
//	result := opera.Async(thunk).Yield()
type Thunk[T any] = func(context.Context) Result[T]

// SuspenseN is a type alias to represent suspended `Thunk` with N-arity arguments.
// A `SuspenseN` yields a `Thunk` with the same parameters as the original function which can be executed by `opera.Async`.
//
// Example:
//
// ctx := context.Background()
//
//	normalOp := func(a int, b string) opera.Result[string] {
//	    return opera.Ok(fmt.Sprintf("%d - %s", a, b))
//	}
//
// suspendedOp := opera.Suspend2(normalOp)
// thunk := suspendedOp(10, "hello")
// ch := opera.Async(ctx, thunk)
// result := opera.Await(ch).Yield()
type Suspense0[T any] = func() Thunk[T]
type Suspendable0[T any] = func(context.Context) Result[T]

type Suspense[T, A any] = func(A) Thunk[T]
type Suspendable[T, A any] = func(context.Context, A) Result[T]

type Suspense2[T, A1, A2 any] = func(A1, A2) Thunk[T]
type Suspendable2[T, A1, A2 any] = func(context.Context, A1, A2) Result[T]

type Suspense3[T, A1, A2, A3 any] = func(A1, A2, A3) Thunk[T]
type Suspendable3[T, A1, A2, A3 any] = func(context.Context, A1, A2, A3) Result[T]

type Suspense4[T, A1, A2, A3, A4 any] = func(A1, A2, A3, A4) Thunk[T]
type Suspendable4[T, A1, A2, A3, A4 any] = func(context.Context, A1, A2, A3, A4) Result[T]

type Suspense5[T, A1, A2, A3, A4, A5 any] = func(A1, A2, A3, A4, A5) Thunk[T]
type Suspendable5[T, A1, A2, A3, A4, A5 any] = func(context.Context, A1, A2, A3, A4, A5) Result[T]

type Suspense6[T, A1, A2, A3, A4, A5, A6 any] = func(A1, A2, A3, A4, A5, A6) Thunk[T]

type Suspendable6[T, A1, A2, A3, A4, A5, A6 any] = func(context.Context, A1, A2, A3, A4, A5, A6) Result[T]

type Suspense7[T, A1, A2, A3, A4, A5, A6, A7 any] = func(A1, A2, A3, A4, A5, A6, A7) Thunk[T]

type Suspendable7[T, A1, A2, A3, A4, A5, A6, A7 any] = func(context.Context, A1, A2, A3, A4, A5, A6, A7) Result[T]

type Suspense8[T, A1, A2, A3, A4, A5, A6, A7, A8 any] = func(A1, A2, A3, A4, A5, A6, A7, A8) Thunk[T]

type Suspendable8[T, A1, A2, A3, A4, A5, A6, A7, A8 any] = func(context.Context, A1, A2, A3, A4, A5, A6, A7, A8) Result[T]

type Suspense9[T, A1, A2, A3, A4, A5, A6, A7, A8, A9 any] = func(A1, A2, A3, A4, A5, A6, A7, A8, A9) Thunk[T]

type Suspendable9[T, A1, A2, A3, A4, A5, A6, A7, A8, A9 any] = func(context.Context, A1, A2, A3, A4, A5, A6, A7, A8, A9) Result[T]

// SuspendN creates a new function which delays the execution of Result-returning function represented by `SuspenseN`.
// A `SuspenseN` yields a `Thunk` with the same parameters as the original function which can be executed by `opera.Async`.
//
// Example:
//
// ctx := context.Background()
//
//	normalOp := func(a int, b string) opera.Result[string] {
//	    return opera.Ok(fmt.Sprintf("%d - %s", a, b))
//	}
//
// suspendedOp := opera.Suspend2(normalOp)
// thunk := suspendedOp(10, "hello")
// ch := opera.Async(ctx, thunk)
// result := opera.Await(ch).Yield()
func Suspend0[T any](fn Suspendable0[T]) Suspense0[T] {
	return func() Thunk[T] {
		return func(ctx context.Context) Result[T] {
			return fn(ctx)
		}
	}
}

func Suspend[T, A any](fn Suspendable[T, A]) Suspense[T, A] {
	return func(a A) Thunk[T] {
		return func(ctx context.Context) Result[T] {
			return fn(ctx, a)
		}
	}
}

func Suspend2[T, A1, A2 any](fn Suspendable2[T, A1, A2]) Suspense2[T, A1, A2] {
	return func(a1 A1, a2 A2) Thunk[T] {
		return func(ctx context.Context) Result[T] {
			return fn(ctx, a1, a2)
		}
	}
}

func Suspend3[T, A1, A2, A3 any](fn Suspendable3[T, A1, A2, A3]) Suspense3[T, A1, A2, A3] {
	return func(a1 A1, a2 A2, a3 A3) Thunk[T] {
		return func(ctx context.Context) Result[T] {
			return fn(ctx, a1, a2, a3)
		}
	}
}

func Suspend4[T, A1, A2, A3, A4 any](
	fn Suspendable4[T, A1, A2, A3, A4],
) Suspense4[T, A1, A2, A3, A4] {
	return func(a1 A1, a2 A2, a3 A3, a4 A4) Thunk[T] {
		return func(ctx context.Context) Result[T] {
			return fn(ctx, a1, a2, a3, a4)
		}
	}
}

func Suspend5[T, A1, A2, A3, A4, A5 any](
	fn Suspendable5[T, A1, A2, A3, A4, A5],
) Suspense5[T, A1, A2, A3, A4, A5] {
	return func(a1 A1, a2 A2, a3 A3, a4 A4, a5 A5) Thunk[T] {
		return func(ctx context.Context) Result[T] {
			return fn(ctx, a1, a2, a3, a4, a5)
		}
	}
}

func Suspend6[T, A1, A2, A3, A4, A5, A6 any](
	fn Suspendable6[T, A1, A2, A3, A4, A5, A6],
) Suspense6[T, A1, A2, A3, A4, A5, A6] {
	return func(a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6) Thunk[T] {
		return func(ctx context.Context) Result[T] {
			return fn(ctx, a1, a2, a3, a4, a5, a6)
		}
	}
}

func Suspend7[T, A1, A2, A3, A4, A5, A6, A7 any](
	fn Suspendable7[T, A1, A2, A3, A4, A5, A6, A7],
) Suspense7[T, A1, A2, A3, A4, A5, A6, A7] {
	return func(a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7) Thunk[T] {
		return func(ctx context.Context) Result[T] {
			return fn(ctx, a1, a2, a3, a4, a5, a6, a7)
		}
	}
}

func Suspend8[T, A1, A2, A3, A4, A5, A6, A7, A8 any](
	fn Suspendable8[T, A1, A2, A3, A4, A5, A6, A7, A8],
) Suspense8[T, A1, A2, A3, A4, A5, A6, A7, A8] {
	return func(a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8) Thunk[T] {
		return func(ctx context.Context) Result[T] {
			return fn(ctx, a1, a2, a3, a4, a5, a6, a7, a8)
		}
	}
}

func Suspend9[T, A1, A2, A3, A4, A5, A6, A7, A8, A9 any](
	fn Suspendable9[T, A1, A2, A3, A4, A5, A6, A7, A8, A9],
) Suspense9[T, A1, A2, A3, A4, A5, A6, A7, A8, A9] {
	return func(a1 A1, a2 A2, a3 A3, a4 A4, a5 A5, a6 A6, a7 A7, a8 A8, a9 A9) Thunk[T] {
		return func(ctx context.Context) Result[T] {
			return fn(ctx, a1, a2, a3, a4, a5, a6, a7, a8, a9)
		}
	}
}
