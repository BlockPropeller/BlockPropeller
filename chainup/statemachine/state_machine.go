package statemachine

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

var (
	// ErrInvalidStep is returned when a resource is passed to the
	// state machine, and the state machine does not have the step
	// registered at the provided state.
	ErrInvalidStep = errors.New("invalid resource step")
)

// State definition for a resource being inserted into the state machine.
//
// An individual state can be configured as a terminating state that is
// either successful or failed and whether the state can be repeated in a
// single step.
type State struct {
	Name string

	IsRepeatable bool

	IsFinished   bool
	IsSuccessful bool
}

// NewState returns a new named state instance.
func NewState(name string) State {
	return State{Name: name}
}

// Repeatable tells the state machine that this state can be repeated inside
// a single step.
//
// This is useful for steps that wait for some condition to be
// met before proceeding to the next state.
func (s State) Repeatable() State {
	s.IsRepeatable = true

	return s
}

// Successful tells the state machine that this state serves as one of the terminating
// states for a given resource. The resource finishes the process in a successful state.
//
// Usually a "success" state will be defined to signify a successful operation on a resource.
func (s State) Successful() State {
	s.IsFinished = true
	s.IsSuccessful = true

	return s
}

// Failure tells the state machine that this state serves as one of the terminating
// states for a given resource. The resource finishes the process in a failed state.
//
// Typical states of this type are "cancelled" and "failed".
func (s State) Failure() State {
	s.IsFinished = true
	s.IsSuccessful = false

	return s
}

// IsEqual compares two state objects and returns whether they are of the same type.
func (s State) IsEqual(other State) bool {
	return s.Name == other.Name
}

// IsIn returns true if the called state is contained in the provided state array.
//
// This method is useful for checking if a particular state is in an array of valid
// states for a particular resource.
func (s State) IsIn(states []State) bool {
	for _, other := range states {
		if s.IsEqual(other) {
			return true
		}
	}

	return false
}

// String satisfies the Stringer interface.
func (s State) String() string {
	return s.Name
}

// MarshalJSON satisfies the json.Marshaler interface.
//
// We want to marshal just the state name.
func (s State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Name)
}

// StatefulResource is any resource provided to the state machine and is able to
// retrieve and update its state. That state is then used to find an appropriate
// step to execute inside the state machine.
type StatefulResource interface {
	GetState() State
	SetState(state State)
}

// Resource is a helper struct intended to be embedded into structures that are
// to be used as StatefulResources in a state machine.
//
// If this struct does not satisfy all your requirements, any valid implementation
// of the StatefulResource can be used instead. Resource struct can be copied as a
// starting point it modifications are needed.
type Resource struct {
	State State `json:"state"`
}

// NewResource returns an initialized Resource with an initial state.
func NewResource(initialState State) Resource {
	return Resource{
		State: initialState,
	}
}

// GetState satisfies the StatefulResource interface.
func (res *Resource) GetState() State {
	return res.State
}

// SetState satisfies the StatefulResource interface.
func (res *Resource) SetState(state State) {
	res.State = state
}

// Step represents a single unit of execution inside a state machine.
//
// Each step is associated with a single state which triggers it.
type Step interface {
	Step(ctx context.Context, res StatefulResource) error
}

// StepFn is a helper function for defining steps that don't have any
// dependencies and can be provided as simple methods.
type StepFn func(ctx context.Context, res StatefulResource) error

// Step satisfies the Step interface.
func (fn StepFn) Step(ctx context.Context, res StatefulResource) error {
	return fn(ctx, res)
}

// MachineBuilder is a structure used for constructing an instance of the StateMachine.
type MachineBuilder struct {
	steps       map[string]Step
	validStates []State
}

// Builder initializes the State Machine builder.
func Builder(validStates []State) *MachineBuilder {
	return &MachineBuilder{
		steps:       make(map[string]Step),
		validStates: validStates,
	}
}

// StepFn adds a State and StepFn tuple to the StateMachine step registry.
func (b *MachineBuilder) StepFn(state State, fn StepFn) *MachineBuilder {
	return b.Step(state, fn)
}

// Step adds a State and Step tuple to the StateMachine step registry.
func (b *MachineBuilder) Step(state State, step Step) *MachineBuilder {
	if _, ok := b.steps[state.Name]; ok {
		panic(errors.Errorf("duplicate step for state: %s", state.Name))
	}

	if !state.IsIn(b.validStates) {
		panic(errors.Errorf("invalid state: %s", state.Name))
	}

	b.steps[state.Name] = step

	return b
}

// Build constructs the final StateMachine from builder configuration.
func (b *MachineBuilder) Build() *StateMachine {
	return &StateMachine{
		steps: b.steps,
	}
}

// StateMachine is a generic StateMachine structure that is able to execute over any
// resource that satisfies the StatefulResource interface.
//
// StateMachine has various rules for executing steps, promoting safe usage.
type StateMachine struct {
	steps map[string]Step
}

// Step advances the StateMachine for a single step.
func (sm *StateMachine) Step(ctx context.Context, res StatefulResource) error {
	if res.GetState().IsFinished {
		panic(errors.New("step must not be called on a finished resource state"))
	}

	machineState := res.GetState()

	step, ok := sm.steps[machineState.Name]
	if !ok {
		return ErrInvalidStep
	}

	err := step.Step(ctx, res)
	if err != nil {
		return err
	}

	if res.GetState().IsEqual(machineState) && !machineState.IsRepeatable {
		panic(errors.Errorf("expected state change after state: %s", machineState))
	}

	return nil
}
