package infrastructure

import (
	"chainup.dev/chainup/statemachine"
	"chainup.dev/chainup/statemachine/step"
)

var (
	// StateCreated is the starting point for a provisioning job.
	StateCreated = statemachine.NewState("created")

	// StateFailed is the terminating state representing provisioning job failure.
	// @TODO: Add failure message to job somewhere.
	StateFailed = statemachine.NewState("failed").Failure()

	// ValidStates of a provision.Job.
	ValidStates = []statemachine.State{StateCreated, StateFailed}
)

// ConfigureStateMachine returns a preconfigured StateMachine
// for running provisioning jobs.
func ConfigureStateMachine() *statemachine.StateMachine {
	return statemachine.Builder(ValidStates).
		Step(StateCreated, step.Simple(StateFailed)).
		Build()
}
