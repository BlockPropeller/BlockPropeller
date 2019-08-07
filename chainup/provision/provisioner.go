package provision

import (
	"context"

	"chainup.dev/chainup/terraform"
	"chainup.dev/lib/log"
	"github.com/pkg/errors"
)

// Provisioner runs the provisioning process from start to finish.
type Provisioner struct {
	StateMachine *JobStateMachine

	JobRepository JobRepository

	//@TODO: Abstract away terraform from here?
	Terraform *terraform.Terraform
}

// NewProvisioner returns a new Provisioner instance.
func NewProvisioner(stateMachine *JobStateMachine, jobRepository JobRepository, terraform *terraform.Terraform) *Provisioner {
	return &Provisioner{StateMachine: stateMachine, JobRepository: jobRepository, Terraform: terraform}
}

// Provision starts the provisioning process and returns after it is complete.
func (p *Provisioner) Provision(ctx context.Context, job *Job) error {
	err := p.JobRepository.Create(job)
	if err != nil {
		return errors.Wrap(err, "create job")
	}

	//@TODO: Create resource creation request for machines that need to be created and services that need to be running on top.
	return p.StateMachine.StepToCompletion(ctx, job)
}

// Undo provisioned infrastructure based on the terraform Workspace.
func (p *Provisioner) Undo(ctx context.Context, job *Job) error {
	if job.WorkspaceSnapshot == nil {
		return errors.New("missing workspace snapshot")
	}

	workspace, err := terraform.RestoreWorkspace(job.WorkspaceSnapshot)
	if err != nil {
		return errors.Wrap(err, "restore workspace")
	}
	defer func() {
		log.Debug("cleaning up Terraform workspace")
		log.Closer(workspace)
	}()

	err = p.Terraform.Init(workspace)
	if err != nil {
		return errors.Wrap(err, "init workspace")
	}

	err = p.Terraform.Destroy(workspace)
	if err != nil {
		return errors.Wrap(err, "destroy workspace")
	}

	return nil
}
