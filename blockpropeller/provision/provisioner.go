package provision

import (
	"context"

	"blockpropeller.dev/blockpropeller/terraform"
	"github.com/pkg/errors"
)

// Provisioner runs the provisioning process from start to finish.
type Provisioner struct {
	StateMachine *JobStateMachine

	//@TODO: Abstract away terraform from here?
	Terraform *terraform.Terraform

	ServerDestroyer *ServerDestroyer
}

// NewProvisioner returns a new Provisioner instance.
func NewProvisioner(
	stateMachine *JobStateMachine,
	terraform *terraform.Terraform,
	srvDestroyer *ServerDestroyer,
) *Provisioner {
	return &Provisioner{
		StateMachine:    stateMachine,
		Terraform:       terraform,
		ServerDestroyer: srvDestroyer,
	}
}

// Provision starts the provisioning process and returns after it is complete.
func (p *Provisioner) Provision(ctx context.Context, job *Job) error {
	return p.StateMachine.StepToCompletion(ctx, job)
}

// Undo provisioned infrastructure based on the terraform Workspace.
func (p *Provisioner) Undo(ctx context.Context, job *Job) error {
	if job.Server == nil {
		return errors.New("missing server associated with the job")
	}

	return p.ServerDestroyer.Destroy(ctx, job.Server)
}
