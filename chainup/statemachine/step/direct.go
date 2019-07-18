package step

import (
	"context"

	"chainup.dev/chainup/statemachine"
)

// Simple step that transitions a resource to a provided state.
func Simple(next statemachine.State) statemachine.StepFn {
	return func(ctx context.Context, res statemachine.StatefulResource) error {
		res.SetState(next)
		return nil
	}
}
