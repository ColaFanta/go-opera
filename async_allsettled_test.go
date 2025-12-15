package opera

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAwaitAllSettledAllNilHetero(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	var ch1 chan Result[int] = nil
	var ch2 chan Result[string] = nil

	result := AwaitAllSettled[any](ctx, ch1, ch2)

	is.Equal(2, len(result))
	is.Equal(Empty[Result[any]](), result[0])
	is.Equal(Empty[Result[any]](), result[1])
}

func TestAwaitAllSettledAllNilHomo(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	var ch1 chan Result[int] = nil
	var ch2 chan Result[string] = nil

	result := AwaitAllSettled[int](ctx, ch1, ch2)

	is.Equal(2, len(result))
	is.Equal(Empty[Result[int]](), result[0])
	is.Equal(Empty[Result[int]](), result[1])
}

func TestAwaitAllSettledHetero(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int], 1)
	ch2 := make(chan Result[string], 1)
	ch3 := make(chan Result[float64], 1)

	ch1 <- Ok(42)
	ch2 <- Err[string](errors.New("Test Error"))
	ch3 <- Ok(3.14)

	result := AwaitAllSettled[any](ctx, ch1, ch2, ch3)
	is.Equal(3, len(result))

	is.True(result[0].IsOk())
	is.Equal(42, CastResult[int](result[0]).Yield())

	is.True(result[1].IsErr())
	is.EqualError(result[1].Err(), "Test Error")

	is.True(result[2].IsOk())
	is.Equal(3.14, CastResult[float64](result[2]).Yield())
}

func TestAwaitAllSettledHomo(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int], 1)
	ch2 := make(chan Result[int], 1)
	ch3 := make(chan Result[int], 1)

	ch1 <- Ok(10)
	ch2 <- Err[int](errors.New("Test Error"))
	ch3 <- Ok(30)

	result := AwaitAllSettled[int](ctx, ch1, ch2, ch3)
	is.Equal(3, len(result))

	is.True(result[0].IsOk())
	is.Equal(10, result[0].Yield())

	is.True(result[1].IsErr())
	is.EqualError(result[1].Err(), "Test Error")

	is.True(result[2].IsOk())
	is.Equal(30, result[2].Yield())
}

func TestAwaitAllSettledHomoCastError(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int], 1)
	ch2 := make(chan Result[int], 1)
	ch3 := make(chan Result[int], 1)

	ch1 <- Ok(10)
	ch2 <- Err[int](errors.New("Test Error"))
	ch3 <- Ok(30)

	r := AwaitAllSettled[string](ctx, ch1, ch2, ch3)
	is.Equal(3, len(r))

	is.True(r[0].IsErr())
	is.ErrorIs(r[0].Err(), ErrAsyncAwaitedValueCastError)
}

func TestAwaitAllSettledT(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int], 1)
	ch2 := make(chan Result[int], 1)
	ch3 := make(chan Result[int], 1)

	ch1 <- Ok(10)
	ch2 <- Err[int](errors.New("Test Error"))
	ch3 <- Ok(30)

	result := AwaitAllSettledT(ctx, ch1, ch2, ch3)
	is.Equal(3, len(result))

	is.True(result[0].IsOk())
	is.Equal(10, result[0].Yield())

	is.True(result[1].IsErr())
	is.EqualError(result[1].Err(), "Test Error")

	is.True(result[2].IsOk())
	is.Equal(30, result[2].Yield())
}
