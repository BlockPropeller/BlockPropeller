package provision

import (
	"context"

	"chainup.dev/chainup/statemachine"
	"github.com/pkg/errors"
)

// Provisioner controls the process of servicing jobs for
// infrastructure creation and deployment management on the created infrastructure.
//
// Provisioner takes a user job that defines what infrastructure is needed
// and what deployments should be running on it and controls the provisioning
// process until the desired requirements are satisfied.
type Provisioner struct {
	jobRepo JobRepository

	stateMachine *statemachine.StateMachine
}

// NewProvisioner returns a new Provisioner instance.
func NewProvisioner(jobRepo JobRepository, stateMachine *statemachine.StateMachine) *Provisioner {
	return &Provisioner{
		jobRepo:      jobRepo,
		stateMachine: stateMachine,
	}
}

// Start a new provisioning Job, returning after the Request is submitted for processing.
func (p *Provisioner) Start(job *Job) error {
	panic("Implement me!")
}

// WaitFor a job to reach a finished state.
func (p *Provisioner) WaitFor(job *Job) error {
	err := p.jobRepo.Create(job)
	if err != nil {
		return errors.Wrap(err, "create provision job")
	}

	//@TODO: Create a state machine wrapper that runs the job step inside a transaction.

	for !job.State.IsFinished {
		err := p.stateMachine.Step(context.TODO(), job)
		if err != nil {
			return errors.Wrap(err, "run state machine step")
		}
	}

	return nil
}
