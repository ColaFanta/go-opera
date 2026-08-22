package opera

import (
	"context"
	"errors"
	"fmt"
)

var ErrAsyncChannelClosed = errors.New("channel closed")
var ErrAsyncArgumentNotChannelOfResult = errors.New("argument is not a channel of Result")
var ErrAsyncAwaitedValueCastError = errors.New("awaited value cast error")
var ErrAsyncAllOpsFailed = errors.New("all async operations failed")

// Async runs fn(ctx) (T, error) in a goroutine and returns a channel for its Result.
func Async[T any](ctx context.Context, fn Thunk[T]) <-chan Result[T] {
	ch := make(chan Result[T], 1)
	go func() {
		select {
		case ch <- fn(ctx):
		case <-ctx.Done():
			select {
			case ch <- Err[T](fmt.Errorf("%w: Result[%T]", ctx.Err(), Empty[T]())):
			default:
			}
		}
	}()
	return ch
}

// Await blocks receiving a Result[T] or context cancellation.
// If the channel is nil, it returns Ok result with empty T value immediately.
func Await[T any](ctx context.Context, ch <-chan Result[T]) Result[T] {
	if ch == nil {
		return Ok(Empty[T]())
	}
	select {
	case res, ok := <-ch:
		if !ok {
			return Err[T](fmt.Errorf(
				"%w: Result[%T]",
				ErrAsyncChannelClosed,
				Empty[T](),
			))
		}
		return res
	case <-ctx.Done():
		return Err[T](fmt.Errorf("%w: Result[%T]", ctx.Err(), Empty[T]()))
	}
}

// AwaitAll waits for all channels and returns collected values or first error.
// The order of returned values corresponds to the order of input channels.
// If a channel is nil, the corresponding value will be an empty T.
// Use AwaitAllSettled to collect all results including errors.
func AwaitAll[T any](ctx context.Context, chans ...any) Result[[]T] {
	inctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type outcome struct {
		index int
		value any
		err   error
	}
	outcomes := make(chan outcome, len(chans))
	for i, ch := range chans {
		go func(index int, channel any) {
			var outcome outcome
			outcome.index = index
			if isChannelNil[Result[T]](channel) {
				// do nothing
			} else if chc, ok := channel.(<-chan Result[T]); ok {
				// for homogeneous channels <-chan Result[T]
				if res, err := Await(inctx, chc).Get(); err != nil {
					outcome.err = err
				} else {
					outcome.value = res
				}
				// for heterogeneous channels, T must be any here, so that channels in `chans` may be heterogeneous
			} else if res, err := awaitAnyReflect(inctx, channel).Get(); err != nil {
				outcome.err = err
			} else {
				outcome.value = res
			}
			outcomes <- outcome
		}(i, ch)
	}

	resultVals := make([]any, len(chans))
	for range chans {
		select {
		case outcome := <-outcomes:
			if outcome.err != nil {
				return Err[[]T](outcome.err)
			}
			resultVals[outcome.index] = outcome.value
		case <-ctx.Done():
			return Err[[]T](fmt.Errorf("%w: Result[%T]", ctx.Err(), Empty[T]()))
		}
	}

	castedVals := make([]T, len(resultVals))
	for i, result := range resultVals {
		if result == nil {
			continue
		}
		if v, ok := result.(T); ok {
			castedVals[i] = v
			continue
		}
		return Err[[]T](fmt.Errorf(
			"%w [%T]: %T, %d",
			ErrAsyncAwaitedValueCastError,
			Empty[T](),
			result,
			i,
		))
	}

	return Ok(castedVals)
}

// AwaitAllT is a type-safe wrapper for AwaitAll for homogeneous channels.
func AwaitAllT[T any](ctx context.Context, chans ...<-chan Result[T]) Result[[]T] {
	chansany := make([]any, len(chans))
	for i, ch := range chans {
		chansany[i] = ch
	}
	return AwaitAll[T](ctx, chansany...)
}

// AwaitFirst waits for the first successful Result from the provided channels
// errors are ignored unless all operations are failed.
// There may be [randomness](https://go.dev/tour/concurrency/5) if multiple channels are ready simultaneously.
func AwaitFirst[T any](ctx context.Context, chans ...any) Result[Tuple[T, int]] {
	if !isAtLeastOneChannelNotNil[T](chans...) {
		return Ok(Empty[Tuple[T, int]]())
	}

	inctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type outcome struct {
		result Tuple[T, int]
		err    error
	}
	outcomes := make(chan outcome, len(chans))
	activeChannels := 0
	for _, ch := range chans {
		if !isChannelNil[Result[T]](ch) {
			activeChannels++
		}
	}
	for i, ch := range chans {
		go func(index int, channel any) {
			var outcome outcome
			if isChannelNil[Result[T]](channel) {
				return
			} else if chc, ok := channel.(<-chan Result[T]); ok {
				// for homogeneous channels
				if res, err := Await(inctx, chc).Get(); err != nil {
					outcome.err = err
				} else {
					outcome.result = Tup(res, index)
				}
			} else if res, err := awaitAnyReflect(inctx, channel).Get(); err != nil {
				// for heterogeneous channels, T must be any here, so that channels in `chans` may be heterogeneous
				outcome.err = err
			} else {
				if res == nil {
					outcome.result = Tup(Empty[T](), index)
				} else if val, ok := res.(T); ok {
					outcome.result = Tup(val, index)
				} else {
					outcome.err = fmt.Errorf(
						"%w [%T]: %T, %d",
						ErrAsyncAwaitedValueCastError,
						Empty[T](),
						res,
						index,
					)
				}
			}
			outcomes <- outcome
		}(i, ch)
	}

	errs := make([]error, 0, len(chans))
	for range activeChannels {
		select {
		case outcome := <-outcomes:
			if outcome.err == nil {
				return Ok(outcome.result)
			}
			errs = append(errs, outcome.err)
		case <-ctx.Done():
			return Err[Tuple[T, int]](fmt.Errorf("%w: Result[%T]", ctx.Err(), Empty[Tuple[T, int]]()))
		}
	}

	return Err[Tuple[T, int]](errors.Join(ErrAsyncAllOpsFailed, errors.Join(errs...)))
}

// AwaitFirstT is a type-safe wrapper for AwaitFirst for homogeneous channels.
func AwaitFirstT[T any](ctx context.Context, chans ...<-chan Result[T]) Result[Tuple[T, int]] {
	{
		chansany := make([]any, len(chans))
		for i, ch := range chans {
			chansany[i] = ch
		}
		return AwaitFirst[T](ctx, chansany...)
	}
}
