package infrastructure

var (
	// DeploymentStateRequested is the initial state of any Deployment.
	// It is used to signify that the Deployment is waiting to be processed.
	DeploymentStateRequested DeploymentState = "requested"
	// DeploymentStateRunning is the final success state of a deployment.
	DeploymentStateRunning DeploymentState = "running"

	// ValidDeploymentStates that are recognized by ChainUP.
	ValidDeploymentStates = []DeploymentState{DeploymentStateRequested, DeploymentStateRunning}
)

// DeploymentState defines a valid Deployment state.
type DeploymentState string

// NewDeploymentState returns a new DeploymentState instance.
func NewDeploymentState(state string) DeploymentState {
	return DeploymentState(state)
}

// IsValid checks whether the DeploymentState is one of recognized values.
func (state DeploymentState) IsValid() bool {
	for _, valid := range ValidDeploymentStates {
		if state == valid {
			return true
		}
	}

	return false
}

// String satisfies the Stringer interface.
func (state DeploymentState) String() string {
	return string(state)
}
