package opera

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAsyncResultSuccess(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	success0 := func() (string, error) {
		time.Sleep(time.Millisecond)
		return "Async Success", nil
	}

	var success Thunk[string] = func(ctx context.Context) Result[string] {
		return Try(success0())
	}

	ch := Async(ctx, success)

	result := <-ch

	is.Equal("Async Success", result.Yield())

}

func TestAsyncResultError(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	err0 := func() (string, error) {
		time.Sleep(time.Millisecond)
		return "", errors.New("Async Error")
	}

	var err Thunk[string] = func(ctx context.Context) Result[string] {
		return Try(err0())
	}

	ch := Async(ctx, err)

	result := <-ch

	is.EqualError(result.Err(), "Async Error")
}

func TestAsyncTryIfSuccess(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	success0 := func() (string, bool) {
		time.Sleep(time.Millisecond)
		return "Async TryIf Success", true
	}

	var success Thunk[string] = func(ctx context.Context) Result[string] {
		return TryHave(success0())
	}

	ch := Async(ctx, success)

	result := <-ch

	is.Equal("Async TryIf Success", result.Yield())
}

func TestAsyncTryIfError(t *testing.T) {
	is := assert.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err0 := func() (string, bool) {
		time.Sleep(time.Millisecond)
		return "", false
	}

	var err Thunk[string] = func(ctx context.Context) Result[string] {
		return TryHave(err0())
	}

	ch := Async(ctx, err)

	result := <-ch

	is.ErrorIs(result.Err(), ErrNoSuchElement)
}

func TestAsyncAwait(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	success0 := func() (string, error) {
		time.Sleep(time.Millisecond)
		return "Async Await Success", nil
	}

	var success Thunk[string] = func(ctx context.Context) Result[string] {
		return Try(success0())
	}

	ch := Async(ctx, success)

	result := Await(ctx, ch)

	is.Equal("Async Await Success", result.Yield())
}

func TestAsyncAwaitChannelClosed(t *testing.T) {
	is := assert.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan Result[string], 1)
	close(ch)

	result := Await(ctx, ch)

	is.EqualError(result.Err(), "channel closed: Result[string]")
}

func TestAsyncAwaitContextCanceled(t *testing.T) {
	is := assert.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	success0 := func() (string, error) {
		time.Sleep(time.Millisecond)
		return "Async Await Context Canceled", nil
	}

	var success Thunk[string] = func(ctx context.Context) Result[string] {
		return Try(success0())
	}

	ch := Async(ctx, success)

	result := Await(ctx, ch)

	is.EqualError(result.Err(), "context canceled: Result[string]")
}
