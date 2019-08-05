package provision

import (
	"context"

	"chainup.dev/chainup/statemachine"
	"chainup.dev/chainup/terraform"
	"chainup.dev/chainup/terraform/resource/digitalocean"
	"chainup.dev/lib/log"
	"github.com/pkg/errors"
)

var (
	// StateCreated is the starting point for a provisioning job.
	StateCreated = statemachine.NewState("created")

	// StateCompleted is the terminating state representing a successful provisioning job.
	StateCompleted = statemachine.NewState("completed").Successful()

	// StateFailed is the terminating state representing provisioning server failure.
	// @TODO: Add failure message to job somewhere.
	StateFailed = statemachine.NewState("failed").Failure()

	// ValidStates of a provision.Job.
	ValidStates = []statemachine.State{StateCreated, StateCompleted, StateFailed}
)

// JobStateMachine defines the state machine for running provisioning jobs.
type JobStateMachine struct {
	*statemachine.StateMachine
}

// ConfigureJobStateMachine returns a preconfigured StateMachine
// for running provisioning jobs.
func ConfigureJobStateMachine(tfStep *TerraformStep) *JobStateMachine {
	return &JobStateMachine{
		StateMachine: statemachine.Builder(ValidStates).
			Step(StateCreated, tfStep).
			Build(),
	}
}

// TerraformStep creates a plan for creating new infrastructure,
// executes it against the given cloud provider and waits for the
// provisioning to finish.
type TerraformStep struct {
	tf *terraform.Terraform
}

// NewTerraformStep returns a new TerraformStep instance.
func NewTerraformStep(tf *terraform.Terraform) *TerraformStep {
	return &TerraformStep{tf: tf}
}

// Step satisfies the State Machine step interface.
func (step *TerraformStep) Step(ctx context.Context, res statemachine.StatefulResource) error {
	job := res.(*Job)

	providerSettigns := job.ProviderSettings

	server := job.Server
	sshKey := server.SSHKey

	// Prepare workspace in which to execute Terraform plan.
	workspace, err := terraform.NewWorkspace()
	if err != nil {
		return errors.Wrap(err, "create new workspace")
	}

	defer func() {
		log.Debug("cleaning up Terraform workspace")
		log.Closer(workspace)
	}()

	log.Debug("created Terraform workspace", log.Fields{
		"dir": workspace.WorkDir(),
	})

	workspace.Add(digitalocean.NewProvider(providerSettigns.Credentials))

	log.Debug("using provider", log.Fields{
		"type":        providerSettigns.Type,
		"credentials": providerSettigns.Credentials,
	})

	doSSHKey := digitalocean.NewSSHKey(sshKey.Name, sshKey.EncodedPublicKey())
	log.Debug("using ssh key", log.Fields{
		"pub":  sshKey.EncodedPublicKey(),
		"priv": sshKey.EncodedPrivateKey(),
	})

	doDroplet := digitalocean.NewDroplet(
		server.Name,
		"ubuntu-18-04-x64",
		"fra1",
		"s-1vcpu-1gb",
		[]*digitalocean.SSHKey{doSSHKey},
	)

	workspace.AddResource(doSSHKey, doDroplet)

	err = workspace.Flush()
	if err != nil {
		return errors.Wrap(err, "flush workspace")
	}

	//@TODO: Run terraform plan.

	//@TODO: Run terraform execute.

	return nil
}
