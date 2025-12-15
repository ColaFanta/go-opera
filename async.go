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
	resultVals := make([]any, len(chans))
	successCount := 0
	errch := make(chan error, len(chans))
	inctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i, ch := range chans {
		go func() {
			if isChannelNil[Result[T]](ch) {
				// do nothing
			} else if chc, ok := ch.(<-chan Result[T]); ok {
				// for homogeneous channels <-chan Result[T]
				if res, err := Await(inctx, chc).Get(); err != nil {
					errch <- err
				} else {
					resultVals[i] = res
				}
				// for heterogeneous channels, T must be any here, so that channels in `chans` may be heterogeneous
			} else if res, err := awaitAnyReflect(inctx, ch).Get(); err != nil {
				errch <- err
			} else {
				resultVals[i] = res
			}
			successCount++
			if successCount == len(chans) {
				cancel()
			}
		}()
	}

	select {
	case err := <-errch:
		cancel()
		return Err[[]T](err)
	case <-ctx.Done():
		return Err[[]T](fmt.Errorf("%w: Result[%T]", ctx.Err(), Empty[T]()))
	case <-inctx.Done():
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

	resvalchan := make(chan Tuple[T, int], 1)
	errs := make([]error, len(chans))
	errCount := 0
	for i, ch := range chans {
		go func() {
			if isChannelNil[Result[T]](ch) {
				// do nothing
			} else if chc, ok := ch.(<-chan Result[T]); ok {
				// for homogeneous channels
				if res, err := Await(inctx, chc).Get(); err != nil {
					errs[i] = err
				} else {
					resvalchan <- Tup(res, i)
				}
			} else if res, err := awaitAnyReflect(inctx, ch).Get(); err != nil {
				// for heterogeneous channels, T must be any here, so that channels in `chans` may be heterogeneous
				errs[i] = err
			} else {
				if res == nil {
					resvalchan <- Tup(Empty[T](), i)
				} else if val, ok := res.(T); ok {
					resvalchan <- Tup(val, i)
				} else {
					errs[i] = fmt.Errorf(
						"%w [%T]: %T, %d",
						ErrAsyncAwaitedValueCastError,
						Empty[T](),
						res,
						i,
					)
				}
			}
			errCount++
			if errCount == len(chans) {
				cancel()
			}
		}()
	}

	select {
	case res := <-resvalchan:
		cancel()
		return Ok(res)
	case <-ctx.Done():
		return Err[Tuple[T, int]](fmt.Errorf("%w: Result[%T]", ctx.Err(), Empty[Tuple[T, int]]()))
	case <-inctx.Done():
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
