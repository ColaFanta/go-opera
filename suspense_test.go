package opera

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSuspendFunc(t *testing.T) {
	is := assert.New(t)
	ctx := context.Background()

	op0 := func(context.Context) Result[int] {
		return Ok(100)
	}
	suspendedOp0 := Suspend0(op0)
	thunk0 := suspendedOp0()
	result0 := thunk0(ctx).Yield()
	is.Equal(100, result0)

	op2 := func(ctx context.Context, a int, b string) Result[string] {
		return Ok(fmt.Sprintf("%d - %s", a, b))
	}
	suspendedOp2 := Suspend2(op2)
	thunk2 := suspendedOp2(10, "hello")
	result2 := thunk2(ctx).Yield()
	is.Equal("10 - hello", result2)
}

func TestSuspend(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	op0 := func(context.Context) Result[int] {
		time.Sleep(30 * time.Millisecond)
		return Ok(200)
	}
	asyncOp0 := Async(ctx, Suspend0(op0)())

	op1 := func(ctx context.Context, a int, b string) Result[string] {
		time.Sleep(20 * time.Millisecond)
		return Ok(fmt.Sprintf("%d - %s", a, b))
	}
	asyncOp1 := Async(ctx, Suspend2(op1)(20, "world"))

	res0 := Await(ctx, asyncOp0).Yield()
	res1 := Await(ctx, asyncOp1).Yield()

	is.Equal(200, res0)
	is.Equal("20 - world", res1)

}

func TestSuspendScoped(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()

	acquired := false
	released := false

	resourceOp := func(ctx context.Context) Result[string] {
		return Scoped(func() Result[int] {
			time.Sleep(10 * time.Millisecond)
			acquired = true
			return Ok(300)
		},
			func(res int) Result[string] {
				acquired = true
				return Ok("resource")
			},
			func(res int) {
				time.Sleep(20 * time.Millisecond)
				released = true
			},
		)
	}

	asyncResourceOp := Async(ctx, Suspend0(resourceOp)())

	result := Await(ctx, asyncResourceOp).Yield()

	is.True(acquired)
	is.True(released)
	is.Equal("resource", result)
}
