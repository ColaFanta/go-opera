package opera

import "errors"

// Scoped manages a resource within a monadic context, ensuring that the resource is properly released after use.
func Scoped[T, R any](
	acquire func() Result[T],
	use func(T) Result[R],
	release func(T),
) Result[R] {
	acquiredRes, err := acquire().Get()
	defer release(acquiredRes)
	if err != nil {
		return Err[R](err)
	}

	return use(acquiredRes)
}

// Try2 constructs a Result containing two values and an error, and wraps them into Tuple.
func Try2[T0, T1 any](v0 T0, v1 T1, err error) Result[Tuple[T0, T1]] {
	return Try(Tup(v0, v1), err)
}

// Try3 constructs a Result containing three values and an error, and wraps them into Tuple3.
func Try3[T0, T1, T2 any](v0 T0, v1 T1, v2 T2, err error) Result[Tuple3[T0, T1, T2]] {
	return Try(Tup3(v0, v1, v2), err)
}

// Try4 constructs a Result containing four values and an error, and wraps them into Tuple4.
func Try4[T0, T1, T2, T3 any](v0 T0, v1 T1, v2 T2, v3 T3, err error) Result[Tuple4[T0, T1, T2, T3]] {
	return Try(Tup4(v0, v1, v2, v3), err)
}

// Try5 constructs a Result containing five values and an error, and wraps them into Tuple5.
func Try5[T0, T1, T2, T3, T4 any](
	v0 T0,
	v1 T1,
	v2 T2,
	v3 T3,
	v4 T4,
	err error,
) Result[Tuple5[T0, T1, T2, T3, T4]] {
	return Try(Tup5(v0, v1, v2, v3, v4), err)
}

// Try6 constructs a Result containing six values and an error, and wraps them into Tuple6.
func Try6[T0, T1, T2, T3, T4, T5 any](
	v0 T0,
	v1 T1,
	v2 T2,
	v3 T3,
	v4 T4,
	v5 T5,
	err error,
) Result[Tuple6[T0, T1, T2, T3, T4, T5]] {
	return Try(Tup6(v0, v1, v2, v3, v4, v5), err)
}

// Try7 constructs a Result containing seven values and an error, and wraps them into Tuple7.
func Try7[T0, T1, T2, T3, T4, T5, T6 any](
	v0 T0,
	v1 T1,
	v2 T2,
	v3 T3,
	v4 T4,
	v5 T5,
	v6 T6,
	err error,
) Result[Tuple7[T0, T1, T2, T3, T4, T5, T6]] {
	return Try(Tup7(v0, v1, v2, v3, v4, v5, v6), err)
}

// Try8 constructs a Result containing eight values and an error, and wraps them into Tuple8.
func Try8[T0, T1, T2, T3, T4, T5, T6, T7 any](
	v0 T0,
	v1 T1,
	v2 T2,
	v3 T3,
	v4 T4,
	v5 T5,
	v6 T6,
	v7 T7,
	err error,
) Result[Tuple8[T0, T1, T2, T3, T4, T5, T6, T7]] {
	return Try(Tup8(v0, v1, v2, v3, v4, v5, v6, v7), err)
}

// Try9 constructs a Result containing nine values and an error, and wraps them into Tuple9.
func Try9[T0, T1, T2, T3, T4, T5, T6, T7, T8 any](
	v0 T0,
	v1 T1,
	v2 T2,
	v3 T3,
	v4 T4,
	v5 T5,
	v6 T6,
	v7 T7,
	v8 T8,
	err error,
) Result[Tuple9[T0, T1, T2, T3, T4, T5, T6, T7, T8]] {
	return Try(Tup9(v0, v1, v2, v3, v4, v5, v6, v7, v8), err)
}

// TryYield extracts the value from a value and an error, bailing out if the error is non-nil.
// It is intended for use inside Do and is a shortcut for Try followed by Yield.
func TryYield[T any](val T, err error) T {
	return Try(val, err).Yield()
}

// TryYieldUnit extracts the Unit value from an error-only Result, bailing out if the error is non-nil.
// It is intended for use inside Do and is a shortcut for TryPass followed by Yield.
func TryYieldUnit(err error) Unit {
	return TryPass(err).Yield()
}

// ErrNilYieldError is the sentinel panic used by YieldErr when called with a nil error.
var ErrNilYieldError = errors.New("YieldErr called with nil error")

// YieldErr yields an error directly, bailing out if the error is non-nil.
// It is intended for use inside Do. Calling it with a nil error is misuse and panics.
func YieldErr[T any](err error) T {
	if err == nil {
		panic(ErrNilYieldError)
	}
	return Err[T](err).Yield()
}
