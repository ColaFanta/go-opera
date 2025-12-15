package opera

import (
	"context"
	"fmt"
	"reflect"
)

func isChannelNil[T any](ch any) bool {
	if ch == nil {
		return true
	}
	if chc, ok := ch.(<-chan T); ok {
		return chc == nil
	}
	chr := reflect.ValueOf(ch)
	if chr.Kind() != reflect.Chan {
		panic("must be used on channel type")
	}

	return chr.IsNil()
}

func isAtLeastOneChannelNotNil[T any](chans ...any) bool {
	atLeastOneChannel := false
	for _, ch := range chans {
		if !isChannelNil[Result[T]](ch) {
			atLeastOneChannel = true
			break
		}
	}
	return atLeastOneChannel
}

func awaitAnyReflect(ctx context.Context, ch any) Result[any] {
	rv := reflect.ValueOf(ch)
	if rv.Kind() != reflect.Chan {
		return Err[any](ErrAsyncArgumentNotChannelOfResult)
	}

	// Build select cases
	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: rv,
		},
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		},
	}

	// Wait on either channel
	chosen, recv, ok := reflect.Select(cases)

	if !ok {
		return Err[any](ErrAsyncChannelClosed)
	}

	if chosen == 1 {
		return Err[any](fmt.Errorf("%w: Result[%T]", ctx.Err(), Empty[any]()))
	}

	m := recv.MethodByName("Get")
	if !m.IsValid() {
		return Err[any](ErrAsyncArgumentNotChannelOfResult)
	}

	out := m.Call(nil)
	res := out[0].Interface()
	err := out[1].Interface()
	if e, ok := err.(error); ok && e != nil {
		return Err[any](e)
	}

	return Ok(res)
}
