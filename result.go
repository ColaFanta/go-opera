package opera

import (
	"errors"
	"strings"
)

// Result represents a value that may either be a success (Ok) or a failure (Err).
type Result[T any] struct {
	val T
	err error
}

// IsOk returns true if the Result is Ok.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// Ok creates a successful Result.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Yield yields the value from a Result, bail out if there is an error.
// If under `Do` notation, error will be caught.
func (r Result[T]) Yield() T {
	if r.IsErr() {
		panic(errYieldErrorWrapper{err: r.err})
	}
	return r.val
}

// Err extracts the error from a Result, nil if there is no error.
func (r Result[T]) Err() error {
	return r.err
}

// Get extracts the value and error from a Result.
func (r Result[T]) Get() (T, error) {
	return r.val, r.err
}

// Or returns the value if the Result is Ok, otherwise returns the fallback value.
func (r Result[T]) Or(fallback T) T {
	if r.IsOk() {
		return r.val
	}
	return fallback
}

// OrF returns the value if the Result is Ok, otherwise calls the fallback function with the error and returns its result.
func (r Result[T]) OrF(fallbackFn func(error) T) T {
	if r.IsOk() {
		return r.val
	}
	return fallbackFn(r.err)
}

// OrEmpty returns the value if the Result is Ok, otherwise returns the zero value of T.
func (r Result[T]) OrEmpty() T {
	if r.IsOk() {
		return r.val
	}
	return Empty[T]()
}

// OrPanic panics with the original error. Do not confuse with Yield, which is used in Do notation.
func (r Result[T]) OrPanic() T {
	if r.IsErr() {
		panic(r.err)
	}
	return r.val
}

// Map applies a function to the value if the Result is Ok, otherwise returns the Err unchanged.
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.IsErr() {
		return r
	}
	return Ok(fn(r.val))
}

// Chain applies a function that returns a Result to the value if the Result is Ok, otherwise returns the Err unchanged.
func (r Result[T]) Chain(fn func(T) Result[T]) Result[T] {
	if r.IsErr() {
		return r
	}
	return fn(r.val)
}

// Tap executes a function with the value if the Result is Ok, returning the original Result.
func (r Result[T]) Tap(fn func(T)) Result[T] {
	if r.IsOk() {
		fn(r.val)
	}
	return r
}

// ErrAny is a sentinel error value used to match any error.
var ErrAny = errors.New("[opera] any error sentinel")

// MapErrIs tests the error using errors.Is, and if it matches, replaces it with another error.
// opera.ErrAny can be used to match any error.
func (r Result[T]) MapErrIs(pred error, other error) Result[T] {
	if r.IsErr() && (errors.Is(pred, ErrAny) || errors.Is(r.err, pred)) {
		return Try(r.val, other)
	}
	return r
}

// MapErrorString tests if the error message contains a substring, and if it does, replaces it with another error.
func (r Result[T]) MapErrorString(contains string, other error) Result[T] {
	if r.IsErr() && strings.Contains(r.err.Error(), contains) {
		return Try(r.val, other)
	}
	return r
}

// MapErr applies a function to the error if the Result is Err, otherwise returns the Ok unchanged.
func (r Result[T]) MapErr(fn func(error) error) Result[T] {
	if r.IsOk() {
		return r
	}
	err := fn(r.err)
	if err == nil {
		return Ok(r.val)
	}
	return Try(r.val, err) // This ensures in some case `val` somehow holds non-zero value
}

// Catch applies a function that returns a Result to the error if the Result is Err, otherwise returns the Ok unchanged.
func (r Result[T]) Catch(fn func(error) Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}
	nr := fn(r.err)
	if nr.IsOk() {
		return nr
	}
	return Try(r.val, nr.err) // This ensures in some case `val` somehow holds non-zero value
}

// CatchIs tests the error using errors.Is, and if it matches, recovers by applying a handler function to produce a value.
// opera.AnyError can be used to match any error.
func (r Result[T]) CatchIs(pred error, fallback T) Result[T] {
	if r.IsErr() && (errors.Is(pred, ErrAny) || errors.Is(r.err, pred)) {
		return Ok(fallback)
	}
	return r
}

// CatchString tests if the error message contains a substring, and if it does, recovers by providing a fallback value.
func (r Result[T]) CatchString(contains string, fallback T) Result[T] {
	if r.IsErr() && strings.Contains(r.err.Error(), contains) {
		return Ok(fallback)
	}
	return r
}

// TapErr executes a function with the error if the Result is Err, returning the original Result.
func (r Result[T]) TapErr(fn func(error)) Result[T] {
	if r.IsErr() {
		fn(r.err)
	}
	return r
}

// ToPointer converts the Result to a pointer, returning nil if there is an error.
func (r Result[T]) ToPointer() *T {
	if r.IsErr() {
		return nil
	}
	return &r.val
}

// ToOption converts the Result to an Option, returning None if there is an error.
func (r Result[T]) ToOption() Option[T] {
	if r.IsErr() {
		return None[T]()
	}
	return Some(r.val)
}

// ToAny maybe helpful when need to convert Result[T] to Result[any].
func (r Result[T]) ToAny() Result[any] {
	return Try[any](r.Get())
}
