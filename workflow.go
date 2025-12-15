package opera

type workflow struct {
	rollbacks []func(error)
}

// Workflow creates a new workflow instance.
// Rollbacks are registered per step and executed in reverse order if any step fails.
// Example:
//
//	opera.Do(func() {
//	    flow := opera.Workflow()
//	    result1 := opera.Step(flow, step1, rollback1)
//	    result2 := opera.Step(flow, step2, rollback2)
//	    ...
//	})
func Workflow() *workflow {
	return &workflow{
		rollbacks: []func(error){},
	}
}

// Step executes a workflow step with rollback support along with an workflow instance.
// If the step returns an error, all registered rollbacks are executed in reverse order.
// Caveat: `opera.Step` yield value of `Result` directly like `opera.Must`, .
// Example:
//
//	opera.Do(func() {
//	    flow := opera.Workflow()
//	    result1 := opera.Step(flow, step1, rollback1)
//	    result2 := opera.Step(flow, step2, rollback2)
//	    ...
//	})
func Step[T any](flow *workflow, work func() Result[T], rollback func(error)) T {
	flow.rollbacks = append(flow.rollbacks, rollback)

	r := work().TapErr(func(err error) {
		for i := len(flow.rollbacks) - 1; i >= 0; i-- {
			flow.rollbacks[i](err)
		}
	})

	return r.Yield()
}
