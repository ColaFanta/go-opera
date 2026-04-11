package opera

import (
	"errors"
	"fmt"
)

var ErrPredicateFalse = errors.New("predicate returns false")
var ErrInputCastFailed = errors.New("cast input failed")
var ErrResultCastFailed = errors.New("inner value type cast failed")

// Try constructs a Result from a value and an error.
func Try[T any](val T, err error) Result[T] {
	return Result[T]{val: val, err: err}
}

// TryPass constructs a Result with no value from just an error.
func TryPass(err error) Result[Unit] {
	return Result[Unit]{val: U, err: err}
}

// TryHave constructs a Result from a value and a boolean flag.
func TryHave[T any](val T, ok bool) Result[T] {
	var err error
	if !ok {
		err = fmt.Errorf("%w: %T", ErrNoSuchElement, val)
	}
	return Try(val, err)
}

// TryLookUp attempts to retrieve a value from a map by key, returning a Result.
func TryLookUp[K comparable, T any](m map[K]T, key K) Result[T] {
	val, ok := m[key]
	if !ok {
		return Err[T](fmt.Errorf("%w: key '%v' not found", ErrNoSuchElement, key))
	}
	return Ok(val)
}

// TryAt attempts to retrieve an element from a slice by index, returning a Result.
// Arrays can be passed as arr[:].
func TryAt[T any](x []T, index int) Result[T] {
	if index < 0 || index >= len(x) {
		return Err[T](fmt.Errorf("%w: index %d out of bounds", ErrNoSuchElement, index))
	}
	return Ok(x[index])
}

// TryIf constructs a Result with no value from just a boolean flag.
func TryIf(ok bool) Result[Unit] {
	if !ok {
		return Err[Unit](ErrPredicateFalse)
	}

	return Ok(U)
}

// TryCast attempts to cast a value to type T, returning a Result.
func TryCast[T any](val any) Result[T] {
	if casted, ok := val.(T); ok {
		return Ok(casted)
	}

	return Err[T](fmt.Errorf("%w: expect %T, actual %T", ErrInputCastFailed, *new(T), val))
}

// MustCast extracts the value from a successful Result of a type cast, bails out if there is an error.
func MustCast[T any](val any) T {
	return TryCast[T](val).Yield()
}

// Must extracts the value from a Result, bails out if there is an error.
func Must[T any](val T, err error) T {
	return Try(val, err).Yield()
}

// MustPass bails out if there is an error.
func MustPass(err error) {
	TryPass(err).Yield()
}

// MustHave extracts the value from a Result, bail out if the boolean flag is false.
func MustHave[T any](val T, ok bool) T {
	return TryHave(val, ok).Yield()
}

// MustLookUp retrieves a value from a map by key, bails out if the key is not found.
func MustLookUp[K comparable, T any](m map[K]T, key K) T {
	return TryLookUp(m, key).Yield()
}

// MustAt retrieves an element from a slice by index, bails out if the index is out of bounds.
// Arrays can be passed as arr[:].
func MustAt[T any](x []T, index int) T {
	return TryAt(x, index).Yield()
}

// MustTrue bails out if the boolean flag is false.
func MustTrue(ok bool) {
	TryIf(ok).Yield()
}

// Ok constructs a successful Result from a value.
func Ok[T any](val T) Result[T] {
	return Result[T]{val: val}
}

// Err constructs a failed Result from an error.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// / Scoped manages a resource within a monadic context, ensuring that the resource is properly released after use.
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

// CastResult attempts to convert the value inside a Result[T] to type U.
func CastResult[U, T any](res Result[T]) Result[U] {
	{
		if res.IsErr() {
			return Err[U](res.Err())
		}
		converted, ok := any(res.Yield()).(U)
		if !ok {
			return Err[U](
				fmt.Errorf(
					"%w: expect %T, actual %T",
					ErrResultCastFailed,
					converted,
					res.Yield(),
				),
			)
		}
		return Ok(converted)
	}
}
