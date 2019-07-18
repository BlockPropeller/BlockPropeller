package statemachine_test

import (
	"context"
	"testing"

	"chainup.dev/chainup/statemachine"
	"chainup.dev/lib/test"
	"github.com/pkg/errors"
)

var (
	StateCreated     = statemachine.NewState("created")
	StateSuccess     = statemachine.NewState("success").Successful()
	StateRepeateable = statemachine.NewState("repeatable").Repeatable()
	StateFirstPart   = statemachine.NewState("first_part")
	StateSecondPart  = statemachine.NewState("second_part")
	StateFailure     = statemachine.NewState("failure").Failure()
	StateCancelled   = statemachine.NewState("cancelled").Failure()

	validStates = []statemachine.State{StateCreated, StateRepeateable, StateFirstPart, StateSecondPart, StateSuccess, StateFailure, StateCancelled}
)

type Job struct {
	statemachine.Resource

	Acc int
}

func NewJob() *Job {
	return &Job{
		Resource: statemachine.NewResource(StateCreated),
	}
}

var AddStep = func(amount int, next statemachine.State) statemachine.StepFn {
	return func(ctx context.Context, res statemachine.StatefulResource) error {
		job := res.(*Job)

		job.Acc += amount

		job.SetState(next)

		return nil
	}
}

type MultiplyStep struct {
	Multiplier int
	Next       statemachine.State
}

func NewMultiplyStep(multiplier int, next statemachine.State) *MultiplyStep {
	return &MultiplyStep{Multiplier: multiplier, Next: next}
}

func (ms MultiplyStep) Step(ctx context.Context, res statemachine.StatefulResource) error {
	job := res.(*Job)

	job.Acc *= ms.Multiplier
	job.SetState(ms.Next)

	return nil
}

func TestSimpleStateMachine(t *testing.T) {
	sm := statemachine.Builder(validStates).
		StepFn(StateCreated, AddStep(10, StateSuccess)).
		Build()

	job := NewJob()

	err := runToCompletion(sm, job)

	test.CheckErr(t, "run state machine", err)
	test.AssertIntsEqual(t, "job updated", job.Acc, 10)
}

func TestMultipleSteps(t *testing.T) {
	sm := statemachine.Builder(validStates).
		StepFn(StateCreated, AddStep(10, StateFirstPart)).
		Step(StateFirstPart, NewMultiplyStep(5, StateSecondPart)).
		StepFn(StateSecondPart, AddStep(5, StateSuccess)).
		Build()

	job := NewJob()

	err := runToCompletion(sm, job)

	test.CheckErr(t, "run state machine", err)
	test.AssertIntsEqual(t, "job updated", job.Acc, 55)
}

func TestRepeatableSteps(t *testing.T) {
	//@TODO: TestRepeatableSteps
}

func TestSuccessfulInvocation(t *testing.T) {
	//@TODO: TestSuccessfulInvocation
}

func TestFailingInvocation(t *testing.T) {
	//@TODO: TestFailingInvocation
}

func TestContextualCancellation(t *testing.T) {
	//@TODO: TestContextualCancellation
}

func runToCompletion(sm *statemachine.StateMachine, res statemachine.StatefulResource) error {
	return runToCompletionCtx(context.Background(), sm, res)
}

func runToCompletionCtx(ctx context.Context, sm *statemachine.StateMachine, res statemachine.StatefulResource) error {
	for !res.GetState().IsFinished {
		err := sm.Step(ctx, res)
		if err != nil {
			return errors.Wrap(err, "execute state machine to completion")
		}
	}

	return nil
}
