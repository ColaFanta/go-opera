package opera

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkFlowStep(t *testing.T) {
	is := assert.New(t)

	callOrder := 0
	finalizer1CalledAt := 0
	finalizer2CalledAt := 0

	flow := Workflow()

	step1 := func() Result[int] {
		return Ok(42)
	}
	finalizer1 := func(err error) Result[Unit] {
		callOrder++
		finalizer1CalledAt = callOrder
		return Ok(U)
	}

	step2 := func() Result[string] {
		return Err[string](assert.AnError)
	}
	finalizer2 := func(err error) Result[Unit] {
		callOrder++
		finalizer2CalledAt = callOrder
		return Ok(U)
	}

	r1 := Step(flow, step1, finalizer1)
	is.Equal(42, r1)

	is.Zero(finalizer1CalledAt)
	is.Zero(finalizer2CalledAt)
	is.False(flow.isFinalized)
	is.Panics(assert.PanicTestFunc(func() {
		Step(flow, step2, finalizer2)
	}))
	is.Equal(1, finalizer2CalledAt)
	is.Equal(2, finalizer1CalledAt)
	is.True(flow.isFinalized)

}

func TestWorkFlowAsyncStep(t *testing.T) {
	is := assert.New(t)
	ctx := t.Context()
	callOrder := 0
	finalizer1CalledAt := 0
	finalizer2CalledAt := 0

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
	finalizer1 := func(err error) Result[Unit] {
		callOrder++
		finalizer1CalledAt = callOrder
		return Ok(U)
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
	finalizer2 := func(err error) Result[Unit] {
		callOrder++
		finalizer2CalledAt = callOrder
		return Ok(U)
	}

	r1 := Step(flow, func() Result[int] {
		return step1(ctx)
	}, finalizer1)
	is.Equal(42, r1)
	is.Zero(finalizer1CalledAt)
	is.Zero(finalizer2CalledAt)
	is.False(flow.isFinalized)

	is.Panics(assert.PanicTestFunc(func() {
		Step(flow, func() Result[string] {
			return step2(ctx)
		}, finalizer2)
	}))
	is.Equal(1, finalizer2CalledAt)
	is.Equal(2, finalizer1CalledAt)
	is.True(flow.isFinalized)
}

func TestWorkFlowFinalize(t *testing.T) {
	is := assert.New(t)

	callOrder := 0
	finalizer1CalledAt := 0
	finalizer2CalledAt := 0
	finalizer1Err := assert.AnError
	finalizer2Err := assert.AnError

	flow := Workflow()

	step1 := func() Result[int] {
		return Ok(42)
	}
	finalizer1 := func(err error) Result[Unit] {
		callOrder++
		finalizer1CalledAt = callOrder
		finalizer1Err = err
		return Ok(U)
	}

	step2 := func() Result[string] {
		return Ok("ok")
	}
	finalizer2 := func(err error) Result[Unit] {
		callOrder++
		finalizer2CalledAt = callOrder
		finalizer2Err = err
		return Ok(U)
	}

	r1 := Step(flow, step1, finalizer1)
	is.Equal(42, r1)
	r2 := Step(flow, step2, finalizer2)
	is.Equal("ok", r2)
	is.False(flow.isFinalized)
	is.Zero(finalizer1CalledAt)
	is.Zero(finalizer2CalledAt)

	Finalize(flow)

	is.True(flow.isFinalized)
	is.Equal(1, finalizer2CalledAt)
	is.Equal(2, finalizer1CalledAt)
	is.NoError(finalizer1Err)
	is.NoError(finalizer2Err)
}

func TestNestedWorkflow(t *testing.T) {
	is := assert.New(t)

	flow := Workflow()

	callOrder := 0
	rootflowFinalizer1CalledAt := 0
	rootflowFinalizer2CalledAt := 0
	nestedflowFinalizer1CalledAt := 0
	nestedflowFinalizer2CalledAt := 0

	nestedWork := func() Result[int] {
		return Do(func() int {
			flow := Workflow()

			step1 := func() Result[int] {
				return Ok(100)
			}
			finalizer1 := func(err error) Result[Unit] {
				callOrder++
				nestedflowFinalizer1CalledAt = callOrder
				return Ok(U)
			}

			step2 := func() Result[string] {
				return Err[string](assert.AnError)
			}
			finalizer2 := func(err error) Result[Unit] {
				callOrder++
				nestedflowFinalizer2CalledAt = callOrder
				return Ok(U)
			}

			r1 := Step(flow, step1, finalizer1)
			is.Equal(100, r1)

			Step(flow, step2, finalizer2)
			return r1
		})
	}

	step1 := func() Result[int] {
		return Ok(42)
	}
	finalizer1 := func(err error) Result[Unit] {
		callOrder++
		rootflowFinalizer1CalledAt = callOrder
		return Ok(U)
	}
	finalizer2 := func(err error) Result[Unit] {
		callOrder++
		rootflowFinalizer2CalledAt = callOrder
		return Ok(U)
	}

	r1 := Step(flow, step1, finalizer1)
	is.Equal(42, r1)

	is.Zero(rootflowFinalizer1CalledAt)
	is.Zero(rootflowFinalizer2CalledAt)
	is.Zero(nestedflowFinalizer1CalledAt)
	is.Zero(nestedflowFinalizer2CalledAt)
	is.False(flow.isFinalized)
	is.Panics(assert.PanicTestFunc(func() {
		Step(flow, nestedWork, finalizer2)
	}))
	is.Equal(1, nestedflowFinalizer2CalledAt)
	is.Equal(2, nestedflowFinalizer1CalledAt)
	is.Equal(3, rootflowFinalizer2CalledAt)
	is.Equal(4, rootflowFinalizer1CalledAt)
	is.True(flow.isFinalized)
}

func TestWorkflowPanicsWhenRunningFinalizedWorkflow(t *testing.T) {
	t.Run("after explicit finalize", func(t *testing.T) {
		is := assert.New(t)

		flow := Workflow()
		Step(flow, func() Result[int] {
			return Ok(1)
		}, func(err error) Result[Unit] {
			return Ok(U)
		})

		Finalize(flow)
		is.True(flow.isFinalized)
		is.PanicsWithValue("workflow has already been finalized", func() {
			Step(flow, func() Result[int] {
				return Ok(2)
			}, func(err error) Result[Unit] {
				return Ok(U)
			})
		})
	})

	t.Run("after step error", func(t *testing.T) {
		is := assert.New(t)

		flow := Workflow()
		Step(flow, func() Result[int] {
			return Ok(1)
		}, func(err error) Result[Unit] {
			return Ok(U)
		})

		is.Panics(func() {
			Step(flow, func() Result[int] {
				return Err[int](assert.AnError)
			}, func(err error) Result[Unit] {
				return Ok(U)
			})
		})
		is.True(flow.isFinalized)
		is.PanicsWithValue("workflow has already been finalized", func() {
			Step(flow, func() Result[int] {
				return Ok(2)
			}, func(err error) Result[Unit] {
				return Ok(U)
			})
		})
	})
}
