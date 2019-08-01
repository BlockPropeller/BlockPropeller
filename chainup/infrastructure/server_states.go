package infrastructure

import (
	"chainup.dev/chainup/statemachine"
	"chainup.dev/chainup/statemachine/step"
)

var (
	// StateCreated is the starting point for a provisioning job.
	StateCreated = statemachine.NewState("created")

	// StateCompleted is the terminating state representing a successful server provisioning.
	StateCompleted = statemachine.NewState("completed").Successful()

	// StateFailed is the terminating state representing provisioning server failure.
	// @TODO: Add failure message to job somewhere.
	StateFailed = statemachine.NewState("failed").Failure()

	// ValidStates of a provision.Job.
	ValidStates = []statemachine.State{StateCreated, StateCompleted, StateFailed}
)

// ConfigureStateMachine returns a preconfigured StateMachine
// for running provisioning jobs.
func ConfigureStateMachine() *statemachine.StateMachine {
	return statemachine.Builder(ValidStates).
		Step(StateCreated, step.Simple(StateFailed)).
		Build()
}
