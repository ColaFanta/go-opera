package opera

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkFlowStep(t *testing.T) {
	is := assert.New(t)

	rollback1Called := false
	rollback2Called := false

	flow := Workflow()

	step1 := func() Result[int] {
		return Ok(42)
	}
	rollback1 := func(err error) {
		rollback1Called = true
	}

	step2 := func() Result[string] {
		return Err[string](assert.AnError)
	}
	rollback2 := func(err error) {
		rollback2Called = true
	}

	r1 := Step(flow, step1, rollback1)
	is.Equal(42, r1)

	is.False(rollback1Called)
	is.False(rollback2Called)
	is.Panics(assert.PanicTestFunc(func() {
		Step(flow, step2, rollback2)
	}))
	is.True(rollback2Called)
	is.True(rollback1Called)

}

func TestWorkFlowAsyncStep(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()
	rollback1Called := false
	rollback2Called := false

	flow := Workflow()

	step1 := func(context.Context) Result[int] {
		ch := make(chan Unit)
		go func() {
			time.Sleep(10 * time.Millisecond)
			ch <- U
		}()
		<-ch
		return Ok(42)
	}
	rollback1 := func(err error) {
		rollback1Called = true
	}

	step2 := func(context.Context) Result[string] {
		ch := make(chan Unit)
		go func() {
			time.Sleep(5 * time.Millisecond)
			ch <- U
		}()
		<-ch
		return Err[string](assert.AnError)
	}
	rollback2 := func(err error) {
		rollback2Called = true
	}

	r1 := Step(flow, func() Result[int] {
		return step1(ctx)
	}, rollback1)
	is.Equal(42, r1)
	is.False(rollback1Called)
	is.False(rollback2Called)

	is.Panics(assert.PanicTestFunc(func() {
		Step(flow, func() Result[string] {
			return step2(ctx)
		}, rollback2)
	}))
	is.True(rollback2Called)
	is.True(rollback1Called)
}

func TestNestedWorkflow(t *testing.T) {
	is := assert.New(t)

	flow := Workflow()

	rootflowRollback1Called := false
	rootflowRollback2Called := false
	nestedflowRollback1Called := false
	nestedflowRollback2Called := false

	nestedWork := func() Result[int] {
		return Do(func() int {
			flow := Workflow()

			step1 := func() Result[int] {
				return Ok(100)
			}
			rollback1 := func(err error) {
				nestedflowRollback1Called = true
			}

			step2 := func() Result[string] {
				return Err[string](assert.AnError)
			}
			rollback2 := func(err error) {
				nestedflowRollback2Called = true
			}

			r1 := Step(flow, step1, rollback1)
			is.Equal(100, r1)

			Step(flow, step2, rollback2)
			return r1
		})
	}

	step1 := func() Result[int] {
		return Ok(42)
	}
	rollback1 := func(err error) {
		rootflowRollback1Called = true
	}
	rollback2 := func(err error) {
		rootflowRollback2Called = true
	}

	r1 := Step(flow, step1, rollback1)
	is.Equal(42, r1)

	is.False(rootflowRollback1Called)
	is.False(rootflowRollback2Called)
	is.False(nestedflowRollback1Called)
	is.False(nestedflowRollback2Called)
	is.Panics(assert.PanicTestFunc(func() {
		Step(flow, nestedWork, rollback2)
	}))
	is.True(rootflowRollback1Called)
	is.True(rootflowRollback2Called)
	is.True(nestedflowRollback1Called)
	is.True(nestedflowRollback2Called)
}
