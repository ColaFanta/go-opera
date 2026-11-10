package opera

import (
	"errors"
	"fmt"
	"reflect"
)

// Option is a container for an optional value of type T. If value exists, Option is
// of type Some. If the value is absent, Option is of type None and stores T's zero
// value; callers must not treat None as retaining a value.
type Option[T any] struct {
	val    T
	hasVal bool
}

// IsSome is an alias to IsPresent.
func (o Option[T]) IsSome() bool {
	return o.hasVal
}

// IsNone returns false when value is present.
func (o Option[T]) IsNone() bool {
	return !o.hasVal
}

// Get returns the value and its presence. It returns T's zero value for None.
func (o Option[T]) Get() (T, bool) {
	if o.IsNone() {
		return Empty[T](), false
	}
	return o.val, o.hasVal
}

// Yield yields value if present or bail out instead.
// If under `Do` notation, an `ErrOptionNoSuchElement` will be caught.
func (o Option[T]) Yield() T {
	if o.IsNone() {
		panic(errYieldErrorWrapper{err: ErrNoSuchElement})
	}

	return o.val
}

// Or returns value if present or default value.
func (o Option[T]) Or(fallback T) T {
	if o.IsNone() {
		return fallback
	}

	return o.val
}

// OrEmpty returns value if present or empty value.
func (o Option[T]) OrEmpty() T {
	if o.IsNone() {
		return Empty[T]()
	}
	return o.val
}

// OrPanic returns value if present or panics with ErrNoSuchElement.
func (o Option[T]) OrPanic() T {
	if o.IsNone() {
		panic(ErrNoSuchElement)
	}
	return o.val
}

// Map executes the mapper function if value is present or returns None if absent.
func (o Option[T]) Map[M any](fn func(value T) M) Option[M] {
	val, ok := o.Get()
	if ok {
		return Some(fn(val))
	}
	return None[M]()
}

// FlatMap executes the mapper function if value is present or returns None if absent.
func (o Option[T]) FlatMap[M any](mapper func(value T) Option[M]) Option[M] {
	val, ok := o.Get()
	if ok {
		return mapper(val)
	}
	return None[M]()
}

// Tap executes the function if value is present and returns the original Option.
func (o Option[T]) Tap(fn func(value T)) Option[T] {
	if o.IsSome() {
		fn(o.val)
	}

	return o
}

// TapNone executes the function if value is absent and returns the original Option.
func (o Option[T]) TapNone(fn func()) Option[T] {
	if o.IsNone() {
		fn()
	}

	return o
}

// OrElse executes the function if value is absent or returns the original Option if present.
func (o Option[T]) OrElse(fn func() Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}

	return fn()
}

// ToResult converts Option to Result, returning Err if value is absent.
// If the provided error is nil, ErrNoSuchElement will be used.
func (o Option[T]) ToResult(errs ...error) Result[T] {
	if o.IsSome() {
		return Ok(o.val)
	}

	err := errors.Join(errs...)
	if err == nil {
		err = ErrNoSuchElement
	}

	return Err[T](fmt.Errorf("%w: %T", err, Empty[T]()))
}

// Cast attempts to cast the Option's value to a different type, returning None if the cast fails or the value is absent.
func (o Option[T]) Cast[C any]() Option[C] {
	val, hasVal := o.Get()
	if !hasVal {
		return None[C]()
	}
	if c, ok := any(val).(C); !ok {
		return None[C]()
	} else {
		return Some(c)
	}
}

// ToPointer returns value if present or a nil pointer.
func (o Option[T]) ToPointer() *T {
	if o.IsNone() {
		return nil
	}

	return &o.val
}

// IsZero assists `omitzero` tag introduced in Go 1.24
func (o Option[T]) IsZero() bool {
	if o.IsNone() {
		return true
	}

	var v any = o.val
	if v, ok := v.(zeroer); ok {
		return v.IsZero()
	}

	return reflect.ValueOf(o.val).IsZero()
}
