package opera

import "errors"

type workflow struct {
	finalizers  []func(error) Result[Unit]
	isFinalized bool
}

// Workflow creates a new workflow instance.
// Finalizers are registered per step and executed in reverse order if any step fails.
// Call Finalize when the workflow completes successfully to run the registered finalizers.
// Example:
//
//	opera.Do(func() {
//	    flow := opera.Workflow()
//	    result1 := opera.Step(flow, step1, finalizer1)
//	    result2 := opera.Step(flow, step2, finalizer2)
//	    ...
//	    opera.Finalize(flow)
//	})
func Workflow() *workflow {
	return &workflow{
		finalizers:  []func(error) Result[Unit]{},
		isFinalized: false,
	}
}

// Step executes a workflow step with finalizer support along with a workflow instance.
// If the step returns an error, all registered finalizers are executed in reverse order.
// Caveat: `opera.Step` yield value of `Result` directly like `opera.Must`, .
// Example:
//
//	opera.Do(func() {
//	    flow := opera.Workflow()
//	    result1 := opera.Step(flow, step1, finalizer1)
//	    result2 := opera.Step(flow, step2, finalizer2)
//	    ...
//	    opera.Finalize(flow)
//	})
func Step[T any](flow *workflow, work func() Result[T], finalizer func(error) Result[Unit]) T {
	if flow.isFinalized {
		panic("workflow has already been finalized")
	}

	flow.finalizers = append(flow.finalizers, finalizer)

	r := work().MapErr(func(err error) error {
		for i := len(flow.finalizers) - 1; i >= 0; i-- {
			flow.finalizers[i](err).TapErr(func(e error) {
				err = errors.Join(err, e)
			})
		}
		flow.isFinalized = true
		return err
	})

	return r.Yield()
}

// Finalize executes all registered finalizers in reverse order.
// Call it after all steps succeed to run success-path finalization exactly once.
func Finalize(flow *workflow) {
	if flow.isFinalized {
		panic("workflow has already been finalized")
	}

	var err error
	for i := len(flow.finalizers) - 1; i >= 0; i-- {
		flow.finalizers[i](err).TapErr(func(e error) {
			err = errors.Join(err, e)
		})
	}
	MustPass(err)
	flow.isFinalized = true
}
