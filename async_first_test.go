package opera

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAwaitFirstHomo(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int])
	ch2 := make(chan Result[int])

	go func() {
		time.Sleep(500 * time.Millisecond)
		ch1 <- Ok(42)
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- Ok(7)
	}()

	result := AwaitFirst[int](ctx, ch1, ch2)

	is.True(result.IsOk())
	is.Equal(7, result.Yield().V0)
	is.Equal(1, result.Yield().V1)
}

func TestAwaitFirstUnit(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[Unit])
	ch2 := make(chan Result[Unit])

	go func() {
		time.Sleep(500 * time.Millisecond)
		ch1 <- Ok(U)
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- Ok(U)
	}()

	result := AwaitFirst[Unit](ctx, ch1, ch2)

	is.True(result.IsOk())
	is.Equal(U, result.Yield().V0)
	is.Equal(1, result.Yield().V1)
}

func TestAwaitFirstHetero(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int])
	ch2 := make(chan Result[string])

	go func() {
		time.Sleep(500 * time.Millisecond)
		ch1 <- Ok(42)
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- Ok("Hello")
	}()

	result := AwaitFirst[any](ctx, ch1, ch2)

	is.True(result.IsOk())
	is.Equal(result.Yield().V0, "Hello")
	is.Equal(result.Yield().V1, 1)
}

func TestAwaitFirstChannelClosed(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int])
	ch2 := make(chan Result[string])

	close(ch1)
	close(ch2)

	result := AwaitFirst[any](ctx, ch1, ch2)

	is.True(result.IsErr())
	is.ErrorIs(result.Err(), ErrAsyncChannelClosed)
	is.ErrorIs(result.Err(), ErrAsyncAllOpsFailed)

}

func TestAwaitFirstContextCanceled(t *testing.T) {
	is := assert.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	ch1 := make(chan Result[int])
	ch2 := make(chan Result[string])

	result := AwaitFirst[any](ctx, ch1, ch2)

	is.True(result.IsErr())
	is.ErrorIs(result.Err(), context.Canceled)
}

func TestAwaitFirstFirstFail(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int])
	ch2 := make(chan Result[string])

	err1 := errors.New("First Channel Error")
	go func() {
		time.Sleep(500 * time.Millisecond)
		ch1 <- Err[int](err1)
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- Ok("second success")
	}()

	result := AwaitFirst[any](ctx, ch1, ch2)

	is.True(result.IsOk())
	is.Equal("second success", result.Yield().V0)
	is.Equal(1, result.Yield().V1)
}

func TestAwaitFirstNonChannelArgErr(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(500 * time.Millisecond)
		ch1 <- 42
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- 7
	}()

	r := AwaitFirst[int](ctx, ch1, ch2)

	is.True(r.IsErr())
	is.ErrorIs(r.Err(), ErrAsyncArgumentNotChannelOfResult)
}

func TestAwaitFirstAllNilChannels(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	var ch1 chan Result[int]
	var ch2 chan Result[string]

	result := AwaitFirst[any](ctx, ch1, ch2)

	is.True(result.IsOk())
	is.Equal(Empty[Tuple[any, int]](), result.Yield())
}

func TestAwaitFirstSomeNilChannels(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int])
	var ch2 chan Result[string] = nil

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch1 <- Ok(100)
	}()

	result := AwaitFirst[any](ctx, ch1, ch2)

	is.True(result.IsOk())
	is.Equal(100, result.Yield().V0)
	is.Equal(0, result.Yield().V1)
}

func TestAwaitFirstAllFailed(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int])
	ch2 := make(chan Result[string])

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch1 <- Err[int](errors.New("first failed"))
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- Err[string](errors.New("second failed"))
	}()

	result := AwaitFirst[any](ctx, ch1, ch2)

	is.True(result.IsErr())
	is.ErrorIs(result.Err(), ErrAsyncAllOpsFailed)
}
