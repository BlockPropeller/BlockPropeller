package provision

import (
	"context"

	"blockpropeller.dev/blockpropeller/terraform"
	"github.com/pkg/errors"
)

// Provisioner runs the provisioning process from start to finish.
type Provisioner struct {
	StateMachine *JobStateMachine

	Scheduler *JobScheduler

	//@TODO: Abstract away terraform from here?
	Terraform *terraform.Terraform

	ServerDestroyer *ServerDestroyer
}

// NewProvisioner returns a new Provisioner instance.
func NewProvisioner(
	stateMachine *JobStateMachine,
	jobScheduler *JobScheduler,
	terraform *terraform.Terraform,
	srvDestroyer *ServerDestroyer,
) *Provisioner {
	return &Provisioner{
		StateMachine:    stateMachine,
		Scheduler:       jobScheduler,
		Terraform:       terraform,
		ServerDestroyer: srvDestroyer,
	}
}

// Provision starts the provisioning process and returns after it is complete.
func (p *Provisioner) Provision(ctx context.Context, jobID JobID) error {
	job, err := p.Scheduler.FindScheduled(ctx, jobID)
	if err != nil {
		return errors.Wrap(err, "find job to provision")
	}

	return p.StateMachine.StepToCompletion(ctx, job)
}

// Undo provisioned infrastructure based on the terraform Workspace.
func (p *Provisioner) Undo(ctx context.Context, job *Job) error {
	if job.Server == nil {
		return errors.New("missing server associated with the job")
	}

	return p.ServerDestroyer.Destroy(ctx, job.Server)
}
