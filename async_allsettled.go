package opera

import (
	"context"
	"fmt"
)

// AwaitAllSettled awaits all channels and returns their results as a slice.
// The order of returned results corresponds to the order of input channels.
// If a channel is nil, the corresponding result will be an Ok with empty T.
// Use AwaitAll to return early on first error.
func AwaitAllSettled[T any](ctx context.Context, chans ...any) []Result[T] {
	results := make([]Option[Result[T]], len(chans))
	completeCount := 0

	inctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for i, ch := range chans {
		go func() {
			if isChannelNil[Result[T]](ch) {
				results[i] = Some(Ok(Empty[T]()))
			} else if chc, ok := ch.(<-chan Result[T]); ok {
				// for homogeneous channels <-chan Result[T]
				results[i] = Some(Await(inctx, chc))
			} else if res, err := awaitAnyReflect(inctx, ch).Get(); err != nil {
				// for heterogeneous channels, T must be any here, so that channels in `chans` may be heterogeneous
				results[i] = Some(Err[T](err))
			} else if res == nil {
				results[i] = Some(Ok(Empty[T]()))
			} else if v, ok := res.(T); ok {
				results[i] = Some(Ok(v))
			} else {
				results[i] = Some(Err[T](fmt.Errorf(
					"%w [%T]: %T, %d",
					ErrAsyncAwaitedValueCastError,
					Empty[T](),
					res,
					i,
				)))
			}
			completeCount++
			if completeCount == len(chans) {
				cancel()
			}
		}()
	}

	ress := make([]Result[T], len(results))
	select {
	case <-ctx.Done():
		for i, res := range results {
			if res.IsNone() {
				ress[i] = Err[T](fmt.Errorf("%w: Result[%T]", ctx.Err(), Empty[T]()))
			} else {
				ress[i] = res.Yield()
			}
		}
		return ress
	case <-inctx.Done():
	}

	for i, res := range results {
		ress[i] = res.Yield()
	}
	return ress
}

// AwaitAllSettledT is a type-safe wrapper for AwaitAllSettled for homogeneous channels.
func AwaitAllSettledT[T any](ctx context.Context, chans ...<-chan Result[T]) []Result[T] {
	chansany := make([]any, len(chans))
	for i, ch := range chans {
		chansany[i] = ch
	}
	return AwaitAllSettled[T](ctx, chansany...)
}
