package opera

import (
	"errors"
	"fmt"
)

var ErrPredicateFalse = errors.New("predicate returns false")
var ErrInputCastFailed = errors.New("cast input failed")
var ErrResultCastFailed = errors.New("inner value type cast failed")

// Try constructs a Result from a value and an error, preserving the value even
// when err is non-nil so callers can inspect partial operation output with Get.
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
//
// Deprecated: Use MayLookUp with ToResult(err) instead.
func TryLookUp[K comparable, T any](m map[K]T, key K) Result[T] {
	val, ok := m[key]
	if !ok {
		return Err[T](fmt.Errorf("%w: key '%v' not found", ErrNoSuchElement, key))
	}
	return Ok(val)
}

// TryAt attempts to retrieve an element from a slice by index, returning a Result.
// Arrays can be passed as arr[:].
//
// Deprecated: Use MayAt with ToResult(err) instead.
func TryAt[T any](x []T, index int) Result[T] {
	if index < 0 || index >= len(x) {
		return Err[T](fmt.Errorf("%w: index %d out of bounds", ErrNoSuchElement, index))
	}
	return Ok(x[index])
}

// TryIf constructs a Result with no value from just a boolean flag.
//
// Deprecated: Use `if` directly, plan to remove in the future.
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
//
// Deprecated: Prefer TryCast with Yield for better readability, plan to remove in the future.
func MustCast[T any](val any) T {
	return TryCast[T](val).Yield()
}

// Must extracts the value from a Result, bails out if there is an error.
//
// Deprecated: Prefer TryYield instead, plan to remove in the future.
func Must[T any](val T, err error) T {
	return Try(val, err).Yield()
}

// MustPass bails out if there is an error.
//
// Deprecated: Prefer TryPass with Yield for better readability, plan to remove in the future.
func MustPass(err error) {
	TryPass(err).Yield()
}

// MustHave extracts the value from a Result, bail out if the boolean flag is false.
//
// Deprecated: Prefer TryHave with Yield for better readability, plan to remove in the future.
func MustHave[T any](val T, ok bool) T {
	return TryHave(val, ok).Yield()
}

// MustLookUp retrieves a value from a map by key, bails out if the key is not found.
//
// Deprecated: Prefer TryLookUp with Yield for better readability, plan to remove in the future.
func MustLookUp[K comparable, T any](m map[K]T, key K) T {
	return TryLookUp(m, key).Yield()
}

// MustAt retrieves an element from a slice by index, bails out if the index is out of bounds.
// Arrays can be passed as arr[:].
//
// Deprecated: Prefer TryAt with Yield for better readability, plan to remove in the future.
func MustAt[T any](x []T, index int) T {
	return TryAt(x, index).Yield()
}

// MustTrue bails out if the boolean flag is false.
//
// Deprecated: Prefer TryIf with Yield for better readability, plan to remove in the future.
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
