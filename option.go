package opera

import (
	"fmt"
	"reflect"
)

// Option is a container for an optional value of type T. If value exists, Option is
// of type Some. If the value is absent, Option is of type None.
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

// Get returns value and presence.
func (o Option[T]) Get() (T, bool) {
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

// Map executes the mapper function if value is present or returns None if absent.
func (o Option[T]) Map(fn func(value T) T) Option[T] {
	if o.IsSome() {
		return Some(fn(o.val))
	}

	return o
}

// Chain executes the mapper function if value is present or returns None if absent.
func (o Option[T]) Chain(mapper func(value T) Option[T]) Option[T] {
	if o.IsSome() {
		return mapper(o.val)
	}

	return o
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
func (o Option[T]) ToResult() Result[T] {
	if o.IsSome() {
		return Ok(o.val)
	}

	return Try(o.val, fmt.Errorf("%w: %T", ErrNoSuchElement, o.val))
}

func (o Option[T]) ToAny() Option[any] {
	val, hasVal := o.Get()
	return Option[any]{val: val, hasVal: hasVal}
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
