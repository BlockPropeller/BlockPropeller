package provision

import (
	"github.com/google/wire"
)

// Set is the Wire provider set for the provisioning package
// that does not depend on any underlying dependencies.
var Set = wire.NewSet(
	NewServerProvisioner,
	NewDeploymentProvisioner,
	NewServerDestroyer,

	NewStepProvisionServer,
	NewStepProvisionDeployment,
	ConfigureJobStateMachine,

	NewJobScheduler,
	NewProvisioner,
)
