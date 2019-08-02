package provision

import (
	"context"

	"github.com/pkg/errors"
)

// Provisioner runs the provisioning process from start to finish.
type Provisioner struct {
	StateMachine *JobStateMachine

	JobRepository JobRepository
}

// NewProvisioner returns a new Provisioner instance.
func NewProvisioner(provisionerStateMachine *JobStateMachine, jobRepo JobRepository) *Provisioner {
	return &Provisioner{StateMachine: provisionerStateMachine, JobRepository: jobRepo}
}

// Provision starts the provisioning process and returns after it is complete.
func (p *Provisioner) Provision(ctx context.Context, job *Job) error {
	err := p.JobRepository.Create(job)
	if err != nil {
		return errors.Wrap(err, "create job")
	}

	//@TODO: Create resource creation request for machines that need to be created and services that need to be running on top.
	//@TODO: Kick-off the provisioning process.
	//@TODO: Wait for the process to complete and return the results to the user.
	return p.StateMachine.StepToCompletion(ctx, job)
}
