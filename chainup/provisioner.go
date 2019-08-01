package chainup

import (
	"context"

	"chainup.dev/chainup/infrastructure"
	"chainup.dev/chainup/statemachine"
)

// Provisioner runs the provisioning process from start to finish.
type Provisioner struct {
	ServerStateMachine *statemachine.StateMachine
}

// NewProvisioner returns a new Provisioner instance.
func NewProvisioner(serverStateMachine *statemachine.StateMachine) *Provisioner {
	return &Provisioner{ServerStateMachine: serverStateMachine}
}

// Provision starts the provisioning process and returns after it is complete.
func (p *Provisioner) Provision(ctx context.Context, server *infrastructure.Server) error {
	//@TODO: Create resource creation request for machines that need to be created and services that need to be running on top.
	//@TODO: Kick-off the provisioning process.
	//@TODO: Wait for the process to complete and return the results to the user.
	return p.ServerStateMachine.StepToCompletion(ctx, server)
}
