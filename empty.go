package opera

import "reflect"

type zeroer interface {
	IsZero() bool
}

// Empty is a helper function that returns the zero value of type T.
func Empty[T any]() T {
	var zero T
	return zero
}

// IsEmpty checks if the given value is the zero value of its type.
// It first tries to handle common types explicitly for performance,
// then falls back to using reflection.
func IsEmpty[T any](val T) bool {
	switch v := any(val).(type) {
	case string:
		if v == "" {
			return true
		}
	case int, int8, int16, int32, int64:
		if v == 0 {
			return true
		}
	case uint, uint8, uint16, uint32, uint64:
		if v == 0 {
			return true
		}
	case float32, float64:
		if v == .0 {
			return true
		}
	case bool:
		if !v {
			return true
		}
	case []any:
		if len(v) == 0 {
			return true
		}
	case map[any]any:
		if len(v) == 0 {
			return true
		}
	}

	if z, ok := any(val).(zeroer); ok {
		if z.IsZero() {
			return true
		}
	}

	isEmpty := reflect.ValueOf(&val).Elem().IsZero()
	return isEmpty
}
