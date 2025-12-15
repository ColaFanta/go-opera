package opera

import (
	"errors"
	"fmt"
)

// Do executes a function within a monadic context, capturing any errors that occur.
// If the function executes successfully, its result is wrapped in a successful Result.
// If the function panics (indicating a failure), the panic is caught and converted into an error Result.
func Do[T any](fn func() T) (result Result[T]) {
	defer func() {
		if r := recover(); r != nil {
			var wrapper errYieldErrorWrapper
			if err, ok := r.(error); !ok {
				panic(r)
			} else if !errors.As(err, &wrapper) {
				panic(err)
			}
			result = Err[T](wrapper.err)
		}
	}()
	return Ok(fn())
}

type errYieldErrorWrapper struct {
	err error
}

func (e errYieldErrorWrapper) Error() string {
	return fmt.Sprintf("yielded error: %s", e.err.Error())
}
