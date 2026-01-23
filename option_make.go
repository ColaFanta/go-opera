package opera

import "errors"

var ErrNoSuchElement = errors.New("no such element")

// Some creates an Option with a value.
func Some[T any](val T) Option[T] {
	return Option[T]{
		hasVal: true,
		val:    val,
	}
}

// None creates an Option without a value.
func None[T any]() Option[T] {
	return Option[T]{
		hasVal: false,
	}
}

// MayHave creates an Option from a value and a boolean flag.
func MayHave[T any](val T, ok bool) Option[T] {
	if ok {
		return Some(val)
	}

	return None[T]()
}

// MayLookUp attempts to retrieve a value from a map by key, returning an Option.
func MayLookUp[T any](m map[string]T, key string) Option[T] {
	val, ok := m[key]
	if !ok {
		return None[T]()
	}

	return Some(val)
}

// MaybeEmpty builds a Some Option when val is not empty, or None.
func MaybeEmpty[T any](val T) Option[T] {
	if IsEmpty(val) {
		return None[T]()
	}

	return Some(val)
}

// MaybeNilPtr builds a Some Option when val is not nil, or None.
func MaybeNilPtr[T any](val *T) Option[T] {
	if val == nil {
		return None[T]()
	}

	return Some(*val)
}

// MayCast tries to cast val to type T and returns Some if successful, or None otherwise.
func MayCast[T any](val any) Option[T] {
	casted, ok := val.(T)
	if !ok {
		return None[T]()
	}

	return Some(casted)
}

// CastOption tries to convert Option[T] to Option[U] and returns None if conversion fails.
func CastOption[U any, T any](opt Option[T]) Option[U] {
	if opt.IsNone() {
		return None[U]()
	}

	converted, ok := any(opt.Yield()).(U)
	if !ok {
		return None[U]()
	}

	return Some(converted)
}
