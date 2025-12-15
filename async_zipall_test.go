package opera

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAsyncZipSomeNil(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int], 1)
	var ch2 chan Result[string]

	ch1 <- Ok(100)

	result := AwaitZip(ctx, ch1, ch2)

	is.True(result.IsOk())
	vals := result.Yield()
	is.Equal(100, vals.V0)
	is.Equal(Empty[string](), vals.V1)
}

func TestAwaitAllHomo(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	tasks := []func() (string, error){
		func() (string, error) {
			time.Sleep(50 * time.Millisecond)
			return "Task 1 Completed", nil
		},
		func() (string, error) {
			time.Sleep(30 * time.Millisecond)
			return "Task 2 Completed", nil
		},
		func() (string, error) {
			time.Sleep(70 * time.Millisecond)
			return "Task 3 Completed", nil
		},
	}

	chs := make([]<-chan Result[string], len(tasks))
	for i, task := range tasks {
		chs[i] = Async(ctx, func(ctx context.Context) Result[string] {
			return Try(task())
		})
	}
	chsany := make([]any, len(chs))
	for i, ch := range chs {
		chsany[i] = ch
	}

	result := AwaitAll[string](ctx, chsany...)

	is.True(result.IsOk())
	is.Equal("Task 2 Completed", result.Yield()[1])
}

func TestAwaitAllHomoError(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	tasks := []func() (string, error){
		func() (string, error) {
			time.Sleep(50 * time.Millisecond)
			return "Task 1 Completed", nil
		},
		func() (string, error) {
			time.Sleep(30 * time.Millisecond)
			return "Task 2 Completed", nil
		},
		func() (string, error) {
			time.Sleep(70 * time.Millisecond)
			return "Task 3 Completed", nil
		},
	}

	chs := make([]<-chan Result[string], len(tasks))
	for i, task := range tasks {
		chs[i] = Async(ctx, func(ctx context.Context) Result[string] {
			return Try(task())
		})
	}
	chsany := make([]any, len(chs))
	for i, ch := range chs {
		chsany[i] = ch
	}

	r := AwaitAll[int](ctx, chsany...)
	is.True(r.IsErr())
	is.ErrorIs(r.Err(), ErrAsyncAwaitedValueCastError)

}

func TestAwaitAllHetero(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	type tmp struct {
		A int
	}

	ch0 := Async(ctx, func(ctx context.Context) Result[int] {
		return Try(1, nil)
	})

	ch1 := Async(ctx, func(ctx context.Context) Result[string] {
		return Try("Success", nil)
	})

	ch2 := Async(ctx, func(ctx context.Context) Result[tmp] {
		return Try(tmp{A: 1}, nil)
	})

	ch3 := Async(ctx, func(ctx context.Context) Result[[]int] {
		return Try([]int{0}, nil)
	})

	result := AwaitZip4(ctx, ch0, ch1, ch2, ch3)

	is.True(result.IsOk())
	val := result.Yield()
	is.Equal(1, val.V0)
	is.Equal("Success", val.V1)
	is.Equal(tmp{A: 1}, val.V2)
	is.Equal([]int{0}, val.V3)

}

func TestAwaitOneFailed(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch0 := Async(ctx, func(ctx context.Context) Result[int] {
		return Try(1, nil)
	})
	ch1 := Async(ctx, func(ctx context.Context) Result[string] {
		return Try("", errors.New("Failed Task"))
	})
	ch2 := Async(ctx, func(ctx context.Context) Result[bool] {
		return Try(true, nil)
	})

	result := AwaitAll[any](ctx, ch0, ch1, ch2)

	is.True(result.IsErr())
	is.EqualError(result.Err(), "Failed Task")
}

func TestAwaitAllNilChannels(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	var ch1 chan Result[int] = nil
	var ch2 chan Result[string] = nil

	result := AwaitAll[any](ctx, ch1, ch2)

	is.True(result.IsOk())
	vals := result.Yield()
	is.Equal(2, len(vals))
	is.Equal(Empty[int](), MayCast[int](vals[0]).OrEmpty())
	is.Equal(Empty[string](), MayCast[string](vals[1]).OrEmpty())
}

func TestAwaitAllSomeNilChannels(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int])
	var ch2 chan Result[string] = nil

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch1 <- Ok(100)
	}()

	result := AwaitAll[any](ctx, ch1, ch2)

	is.True(result.IsOk())
	vals := result.Yield()
	is.Equal(2, len(vals))
	is.Equal(100, MayCast[int](vals[0]).OrEmpty())
	is.Equal(Empty[string](), MayCast[string](vals[1]).OrEmpty())
}

func TestAsyncNonChannelArgError(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- 42
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- 7
	}()

	r := AwaitAll[int](ctx, ch1, ch2)

	is.True(r.IsErr())
	is.ErrorIs(r.Err(), ErrAsyncArgumentNotChannelOfResult)

}

func TestAsyncAwaitT(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int], 1)
	ch2 := make(chan Result[int], 1)
	ch3 := make(chan Result[int], 1)

	ch1 <- Ok(10)
	ch2 <- Ok(20)
	ch3 <- Ok(30)

	result := AwaitAllT(ctx, ch1, ch2, ch3)

	is.True(result.IsOk())
	vals := result.Yield()
	is.Equal([]int{10, 20, 30}, vals)
}

func TestAsyncAwaitTError(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	ch1 := make(chan Result[int], 1)
	ch2 := make(chan Result[int], 1)
	ch3 := make(chan Result[int], 1)

	ch1 <- Ok(10)
	ch2 <- Err[int](errors.New("Test Error"))
	ch3 <- Ok(30)

	result := AwaitAllT(ctx, ch1, ch2, ch3)

	is.True(result.IsErr())
	is.EqualError(result.Err(), "Test Error")
}
